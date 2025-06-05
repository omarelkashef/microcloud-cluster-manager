#!/usr/bin/env python3
# Copyright 2025 Canonical Ltd.
# See LICENSE file for licensing details.

import logging
import ops

logger = logging.getLogger(__name__)


class ClusterManagerCharm(ops.CharmBase):
    def __init__(self, framework: ops.Framework) -> None:
        super().__init__(framework)
        self.pebble_service_name = "microcloud-cluster-manager"
        framework.observe(self.on.microcloud_cluster_manager_pebble_ready, self._on_microcloud_cluster_manager_pebble_ready)


    def _on_microcloud_cluster_manager_pebble_ready(self, event: ops.PebbleReadyEvent) -> None:
        container = event.workload
        container.replan()
        self.unit.status = ops.ActiveStatus()


if __name__ == "__main__":  # pragma: nocover
    ops.main(ClusterManagerCharm)
