package main

import (
	"crypto/md5"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/ncw/swift"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Path struct {
	file_path string
	obj_path  string
	dnl_path  string
}

type Timing struct {
	Upl float64
	Dnl float64
	Del float64
	Len int64
	Obj string
}

func getHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func writeLog(log *os.File, str string) {
	if _, err := log.WriteString(str); err != nil {
		fmt.Println("ERROR: Could not write to log file")
		fmt.Println(err)
	}
}

func main() {
	var err error
	label := flag.String("label", "label", "'snake_case' representation of the provider used in the 'object_load' config.py")
	batch := flag.Int("batch", 1, "The number of times to run the load test")
	dir := flag.String("dir", "uploads/small", "Absolute or relative path to a directory to be uploaded")
	bucket := flag.String("bucket", "global_unique_bucket_name_12345", "The bucket that should be used for testing (auto created/deleted)")
	endpoint := flag.String("endpoint", "https://auth-east.cloud.ca/v2.0", "The swift object storage public url")
	identity := flag.String("identity", "", "Your swift (auth v2.0) object storage identity")
	password := flag.String("password", "", "Your swift (auth v2.0) object storage password")
	insecure := flag.Bool("insecure", false, "Do not verify the SSL connection endpoint")
	flag.Parse()

	if *dir == "" || *bucket == "" || *identity == "" || *password == "" {
		fmt.Println("\nERROR: 'dir', 'bucket', 'identity' and 'password' are required\n")
		flag.Usage()
		os.Exit(2)
	}

	parts := strings.Split(*identity, ":")
	var tenant, username string
	if len(parts) > 1 {
		tenant = parts[0]
		username = parts[1]
	} else {
		fmt.Println("\nERROR: The 'identity' needs to be formated as '<tenant>:<username>'\n")
		flag.Usage()
		os.Exit(2)
	}

	// make dir absolute so it is easier to work with
	abs_dir, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Println("\nERROR: Problem resolving the specified directory\n")
		os.Exit(2)
	}
	_, ctx_dir := filepath.Split(abs_dir)

	// make the downloads directory if it does not exist
	err = os.MkdirAll("downloads", 0777)
	if err != nil {
		fmt.Println("\nERROR: Problem creating 'downloads' directory")
		fmt.Println(err)
		fmt.Println("")
		os.Exit(2)
	}

	// make the logs directory if it does not exist
	err = os.MkdirAll("logs", 0777)
	if err != nil {
		fmt.Println("\nERROR: Problem creating 'logs' directory")
		fmt.Println(err)
		fmt.Println("")
		os.Exit(2)
	}

	// create the log file if it does not exist
	log_file := fmt.Sprintf("logs%s%s_swift_%s.log", string(os.PathSeparator), *label, ctx_dir)
	if _, err := os.Stat(log_file); os.IsNotExist(err) {
		f, err := os.Create(log_file)
		if err != nil {
			fmt.Printf("\nERROR: Problem creating the log file '%s'\n", log_file)
			fmt.Println(err)
		}
		f.Close()
	}
	// open the log file for appending
	log, err := os.OpenFile(log_file, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("\nERROR: Could not create the log file '%s'\n", log_file)
		fmt.Println(err)
		fmt.Println("")
		os.Exit(2)
	}
	defer log.Close()

	// run the whole process 'batch' number of times
	for n := 0; n < *batch; n++ {
		total := 0.0

		writeLog(log, fmt.Sprintf(":using the swift api on %s for directory '%s'\n", *label, ctx_dir))

		transport := &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			MaxIdleConnsPerHost: 2048,
		}
		if *insecure {
			transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}

		// make a swift connection
		conn := swift.Connection{
			Tenant:    tenant,
			UserName:  username,
			ApiKey:    *password,
			AuthUrl:   *endpoint,
			Transport: transport,
		}

		// authenticate swift user
		start := time.Now()
		err = conn.Authenticate()
		if err != nil {
			fmt.Println("\nERROR: Authentication failed.")
			fmt.Println(err)
			fmt.Println("")
			os.Exit(2)
		}
		elapsed := time.Since(start).Seconds()
		total += elapsed
		writeLog(log, fmt.Sprintf("%-12.4f:create connection\n", elapsed))

		// create the container if it does not already exist
		start = time.Now()
		err = conn.ContainerCreate(*bucket, nil)
		if err != nil {
			fmt.Println("\nERROR: Problem creating the specified bucket")
			fmt.Println(err)
			os.Exit(2)
		}
		elapsed = time.Since(start).Seconds()
		total += elapsed
		writeLog(log, fmt.Sprintf("%-12.4f:create container\n", elapsed))
		fmt.Printf("Using bucket: %s\n", *bucket)
		fmt.Println("Starting upload...  This can take a while, go get a coffee.  :)")

		// walk the file system and pull out the important info (because 'Walk' is a blocking function)
		dirs := make([]*Path, 0)
		objs := make([]*Path, 0)
		err = filepath.Walk(abs_dir, func(path string, info os.FileInfo, _ error) (err error) {
			obj_path := strings.TrimPrefix(path, abs_dir)                     // remove abs_dir from path
			obj_path = strings.TrimPrefix(obj_path, string(os.PathSeparator)) // remove leading slash if it exists
			dnl_path := fmt.Sprintf("downloads%s%s", string(os.PathSeparator), obj_path)
			obj_path = filepath.ToSlash(obj_path) // fix windows paths
			if len(obj_path) > 0 {
				if info.IsDir() {
					dirs = append(dirs, &Path{
						obj_path: obj_path,
						dnl_path: dnl_path,
					})
				} else {
					if info.Mode().IsRegular() {
						objs = append(objs, &Path{
							file_path: path,
							obj_path:  obj_path,
							dnl_path:  dnl_path,
						})
					}
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println("\nERROR: Problem discovering a file\n")
			fmt.Println(err)
			os.Exit(2)
		}

		// channels
		resc, errc := make(chan Timing), make(chan error)

		// put all the remote directory objects and local download directories in place
		var dir_wg sync.WaitGroup // use wait group since we are not doing any special handling of goroutine
		for _, p := range dirs {
			dir_wg.Add(1)
			go func(obj_path, dnl_path string) error {
				defer dir_wg.Done()
				// create the object directory
				obj, _, err := conn.Object(*bucket, obj_path)
				if err == nil && obj.ContentType == "application/directory" {
					fmt.Printf("unchanged: %s\n", obj_path)
				} else {
					err = conn.ObjectPutString(*bucket, obj_path, "", "application/directory")
					if err != nil {
						fmt.Printf("\nERROR: Problem creating directory '%s'\n", obj_path)
						fmt.Println(err)
						return err
					}
					fmt.Printf("added object dir: %s\n", obj_path)
				}
				// create the local download directory (if needed) so download does not fail later
				err = os.MkdirAll(dnl_path, 0777)
				if err != nil {
					fmt.Printf("ERROR: Problem creating '%s' directory\n", dnl_path)
					fmt.Println(err)
					return err
				}
				return nil
			}(p.obj_path, p.dnl_path)
		}
		dir_wg.Wait()

		// upload all the objects into their respective remote directories
		for _, p := range objs {
			go func(path, obj_path, dnl_path string) {
				timing := &Timing{0, 0, 0, 0, obj_path}
				hash, err := getHash(path)
				if err != nil {
					fmt.Printf("ERROR: Problem creating object hash\n")
					errc <- err
					return
				}

				// upload
				obj, _, err := conn.Object(*bucket, obj_path)
				if err != nil || obj.Hash != hash {
					fmt.Printf("  started: %s\n", obj_path)
					f, err := os.Open(path)
					if err != nil {
						fmt.Printf("ERROR: Problem opening file '%s'\n", path)
						errc <- err
						return
					}
					defer f.Close()
					start = time.Now()
					_, err = conn.ObjectPut(*bucket, obj_path, f, true, hash, "", nil)
					if err != nil {
						fmt.Printf("ERROR: Problem uploading object '%s'\n", obj_path)
						errc <- err
						return
					}
					timing.Upl = time.Since(start).Seconds()
					fmt.Printf(" uploaded: %s\n", obj_path)
				} else {
					fmt.Printf(" unchanged: %s\n", obj_path)
				}

				// get the size
				obj, _, err = conn.Object(*bucket, obj_path)
				if err != nil {
					fmt.Printf("ERROR: Problem with the uploaded object '%s'\n", obj_path)
					errc <- err
					return
				}
				timing.Len = obj.Bytes

				// download
				f, err := os.Create(dnl_path)
				if err != nil {
					fmt.Printf("ERROR: Problem creating file '%s'\n", dnl_path)
					errc <- err
					return
				}
				defer f.Close()
				start = time.Now()
				_, err = conn.ObjectGet(*bucket, obj_path, f, true, nil)
				if err != nil {
					fmt.Printf("ERROR: Problem downloading object '%s'\n", obj_path)
					errc <- err
					return
				}
				timing.Dnl = time.Since(start).Seconds()
				fmt.Printf(" downloaded: %s\n", obj_path)

				// delete
				start = time.Now()
				err = conn.ObjectDelete(*bucket, obj_path)
				if err != nil {
					fmt.Printf("ERROR: Problem deleting object '%s'\n", obj_path)
					errc <- err
					return
				}
				timing.Del = time.Since(start).Seconds()
				fmt.Printf(" deleted: %s\n", obj_path)
				resc <- *timing
				return
			}(p.file_path, p.obj_path, p.dnl_path)
		}
		for i := 0; i < len(objs); i++ {
			select {
			case res := <-resc:
				writeLog(log, fmt.Sprintf("%-12.4f%-12d:uploading object - %s\n", res.Upl, res.Len, res.Obj))
				writeLog(log, fmt.Sprintf("%-12.4f%-12d:downloading object - %s\n", res.Dnl, res.Len, res.Obj))
				writeLog(log, fmt.Sprintf("%-12.4f%-12d:deleting object - %s\n", res.Del, res.Len, res.Obj))
				total += res.Upl + res.Dnl + res.Del
			case err := <-errc:
				fmt.Println(err)
			}
		}

		// delete any remote directory objects we created so we can delete the container
		var del_wg sync.WaitGroup // use wait group since we are not doing any special handling of goroutine
		for _, p := range dirs {
			del_wg.Add(1)
			go func(obj_path string) error {
				defer del_wg.Done()
				err = conn.ObjectDelete(*bucket, obj_path)
				if err != nil {
					fmt.Printf("ERROR: Problem deleting object '%s'\n", obj_path)
					fmt.Println(err)
					return err
				}
				fmt.Printf("deleted object dir: %s\n", obj_path)
				return nil
			}(p.obj_path)
		}
		del_wg.Wait()

		// delete the container to clean up the load test
		start = time.Now()
		err = conn.ContainerDelete(*bucket)
		if err != nil {
			fmt.Printf("ERROR: Problem deleting bucket '%s'\n", *bucket)
			fmt.Println(err)
		} else {
			fmt.Printf("removed bucket '%s'\n", *bucket)
		}
		elapsed = time.Since(start).Seconds()
		total += elapsed
		writeLog(log, fmt.Sprintf("%-12.4f:delete bucket\n", elapsed))

		// print the total
		writeLog(log, fmt.Sprintf("----------\n%-12.4f:total time for operations\n\n\n", total))

		// clean up the local downloads directory
		os.RemoveAll("downloads")
		fmt.Println("removed local temp 'downloads' directory\n")
	}
	fmt.Println("Load test finished!  Check the 'logs' directory for details...\n")
}
