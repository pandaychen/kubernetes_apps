FROM centos:7
ADD server /root/server
ADD client /root/client
RUN chmod +x /root/server \
	&& chmod +x /root/client
EXPOSE 11111
CMD /root/server -p 11111

