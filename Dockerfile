FROM ubuntu:22.04

COPY boss_apu_exporter /usr/bin/boss_apu_exporter

EXPOSE 9835
ENTRYPOINT ["/usr/bin/boss_apu_exporter"]
