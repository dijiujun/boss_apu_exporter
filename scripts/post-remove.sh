#!/bin/sh
set -e

userdel -f boss_apu_exporter || true

systemctl daemon-reload
