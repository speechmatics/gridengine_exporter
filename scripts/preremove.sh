#!/bin/sh

set -e

systemctl unmask gridengine_exporter
systemctl stop gridengine_exporter
