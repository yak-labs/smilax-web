FROM busybox:latest
ADD _chirp /opt/chirp
ADD smilax5.tcl /opt/smilax5.tcl
ADD _lib_ /lib
RUN chmod +x /opt/* /lib/*
CMD ["/opt/chirp", "/opt/smilax5.tcl", "/mnt"]
