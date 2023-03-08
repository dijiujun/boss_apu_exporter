#!/bin/sh
set -e

systemctl stop boss_apu_exporter.service || true
systemctl disable boss_apu_exporter.service || true

systemctl daemon-reload
