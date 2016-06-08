
swift_load
==========

This is a cross platform tool which generates load on an OpenStack Swift environment to do performance testing.  The tool outputs files into a `logs` directory which can then be used by the [`object_load`](https://github.com/swill/object_load) tool to create graphs (contact wstevens@cloudops.com if interested in [`object_load`](https://github.com/swill/object_load)).


SETUP
-----

The code is not required to use the script.  The executables are available in the `./bin` directory.

USAGE
-----

The script can be run on any directory of files, but the default directory is `uploads/small` as per the [`object_load`](https://github.com/swill/object_load) file generation scripts.  If you would like to test with a different directory, simply pass in the `-dir="abs/or/rel/to/dir` parameter.

The usage documentation for the script is accessible through the `-h` or `-help` flags.

``` bash
$ ./swift_load -h
Usage of ./swift_load:
  -batch int
      The number of times to run the load test (default 1)
  -bucket string
      The bucket that should be used for testing (auto created/deleted) 
      (default "global_unique_bucket_name_12345")
  -dir string
      Absolute or relative path to a directory to be uploaded 
      (default "uploads/small")
  -endpoint string
      The swift object storage public url 
      (default "https://auth-east.cloud.ca/v2.0")
  -identity string
      Your swift (auth v2.0) object storage identity
  -insecure
      Do not verify the SSL connection endpoint
  -label string
      'snake_case' representation of the provider used in the 'object_load' 
      config.py (default "label")
  -password string
      Your swift (auth v2.0) object storage password
```

An example run would look like the following:

```
$ ./swift_load -identity="check_your_profile" -password="check_your_profile"
Using bucket: global_unique_bucket_name_12345
Starting upload...  This can take a while, go get a coffee.  :)
  started: .DS_Store
  started: 100.txt
  started: 200.txt
  started: 600.txt
  started: 150.txt
  started: 400.txt
  started: 250.txt
  started: 700.txt
  started: 50.txt
 uploaded: 150.txt
 uploaded: 200.txt
  started: 1000.txt
  started: 950.txt
  started: 800.txt
  started: 550.txt
  started: 350.txt
  started: 850.txt
  started: 900.txt
  started: 500.txt
  started: 450.txt
  started: 650.txt
 downloaded: 150.txt
  started: 300.txt
 downloaded: 200.txt
  started: 750.txt
 uploaded: .DS_Store
 uploaded: 50.txt
 deleted: 150.txt
 deleted: 200.txt
 uploaded: 400.txt
 downloaded: 50.txt
 downloaded: 400.txt
 deleted: 50.txt
 deleted: 400.txt
 uploaded: 250.txt
 downloaded: 250.txt
 deleted: 250.txt
 uploaded: 350.txt
 uploaded: 300.txt
 downloaded: 350.txt
 downloaded: 300.txt
 deleted: 350.txt
 uploaded: 550.txt
 deleted: 300.txt
 downloaded: 550.txt
 deleted: 550.txt
 uploaded: 500.txt
 downloaded: 500.txt
 deleted: 500.txt
 uploaded: 800.txt
 uploaded: 700.txt
 downloaded: 800.txt
 downloaded: 700.txt
 deleted: 800.txt
 deleted: 700.txt
 uploaded: 650.txt
 downloaded: 650.txt
 deleted: 650.txt
 uploaded: 950.txt
 uploaded: 100.txt
 downloaded: 950.txt
 deleted: 950.txt
 downloaded: .DS_Store
 downloaded: 100.txt
 deleted: .DS_Store
 uploaded: 1000.txt
 deleted: 100.txt
 downloaded: 1000.txt
 deleted: 1000.txt
 uploaded: 450.txt
 uploaded: 600.txt
 downloaded: 450.txt
 downloaded: 600.txt
 deleted: 450.txt
 deleted: 600.txt
 uploaded: 750.txt
 uploaded: 850.txt
 downloaded: 750.txt
 uploaded: 900.txt
 deleted: 750.txt
 downloaded: 850.txt
 downloaded: 900.txt
 deleted: 850.txt
 deleted: 900.txt
removed bucket 'global_unique_bucket_name_12345'
removed local temp 'downloads' directory

Load test finished!  Check the 'logs' directory for details...
```


BUILDING FROM SOURCE
--------------------

If you want to run from source you would do the following.

``` bash
$ git clone https://github.com/swill/swift_load.git
$ cd swift_load
$ go build
$ ./swift_load -h
```


CROSS COMPILING
---------------

To cross compile this code, you can use `gox`.  I have included a basic builder in `./_build.sh` to get you going.

Learn more about installing `gox` at: [https://github.com/mitchellh/gox](https://github.com/mitchellh/gox)

Compilation process:
``` bash
$ cd /path/to/swift_load
$ ./_build.sh
```

