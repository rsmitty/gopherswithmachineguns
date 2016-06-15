##Gophers With Machine Guns

This application is a spin of the very cool [Bees With Machine Guns](https://github.com/newsapps/beeswithmachineguns) project created by the Chicago Tribune. A couple of differences are immediately clear, however. This project is written in golang and targets OpenStack by making use of the Gophercloud tool to create and delete servers.

####Usage
Most Openstack auth info can be set by simply sourcing your keystone rc file. However, this project expects many flags in order to complete the full lifecycle of testing. A sample run would look like:

```bash
go run main.go -count 3 -image "myimage" -flavor "m1.small" -network "private-net" \
-floating-network "public-net" -keyname spencer-key -sshuser "ubuntu" -sshkey "/path/to/id_rsa" \
-endpoint "http://endpoint-to-test.com/" -sim-reqs 250  -tot-reqs 10000
```

All of these flags are more or less required, with the exception of floating-network, which is optional.

####Disclaimer
Similar to the disclaimer on Bees With Machine Guns, this is pretty much an easy way to create a DDOS attack of an endpoint. Make sure you own said endpoint before using this tool.
 
####Future Work
- More detail on usage and flags in the readme
- Learn how to create go binaries
- Output reports from the test. Currently just outputs the `ab` tool's output to stdout.
- General refactoring. As I learn more about golang, I feel quite sure that I'll find some bonehead moves I've made in this project.
