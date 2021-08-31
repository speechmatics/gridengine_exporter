#!/bin/sh

set -e

systemctl daemon-reload
systemctl enable gridengine_exporter
systemctl restart gridengine_exporter
