
swift_load
==========

This is a cross platform tool which generates load on an OpenStack Swift environment to do performance testing.  The tool outputs files into a `logs` directory which can then be used by the `object_load` tool to create graphs (contact wstevens@cloudops.com if interested in `object_load`).

The code is not required to use the script.  The executables are already available on Cloud.ca Object Storage under `Files > tools > swift_load`.


USAGE
-----

The script can be run on any directory of files, but the default directory is `uploads/small` as per the `object_load` file generation scripts.  If you would like to test with a different directory, simply pass in the `-dir="abs/or/rel/to/dir` parameter.

The usage documentation for the script is accessible through the `-h` or `-help` flags.

``` bash
$ ./swift_load -h
Usage of ./swift_load:
  -batch=1: The number of times to run the load test
  -bucket="global_unique_bucket_name_12345": The bucket that should be used for testing (auto created/deleted)
  -dir="uploads/small": Absolute or relative path to a directory to be uploaded
  -endpoint="https://auth-east.cloud.ca/v2.0": The swift object storage public url
  -identity="": Your swift (auth v2.0) object storage identity
  -insecure=false: Do not verify the SSL connection endpoint
  -label="label": 'snake_case' representation of the provider used in the 'object_load' config.py
  -password="": Your swift (auth v2.0) object storage password
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
$ git clone http://git.cloudops.net/eng/swift_load.git
$ cd swift_load
$ go build
$ ./swift_load -h
```


CROSS COMPILING
---------------

Using the script from source is not ideal, instead it should be compiled and the executable should be distributed.  Since this is written in Go (golang), it will have to be compiled for each OS independently.  There is an excellent package called `goxc` which enables you to compile for all OS platforms at the same time.

Learn more about installing `goxc` at: [https://github.com/laher/goxc](https://github.com/laher/goxc)

Compilation process:
``` bash
$ cd /path/to/swift_load
$ goxc
```

