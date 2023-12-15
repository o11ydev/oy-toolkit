---
title: Mimir resources
---

This tool helps estimate the necessary CPU, memory, and disk space for Grafana
Mimir setups. Input details like the number of active series and queries per
second, and the calculator provides resource requirements for each system
component.  It is based on [Grafana's Planning Grafana Mimir capacity
requirements](https://grafana.com/docs/mimir/latest/manage/run-production-environment/planning-capacity).

{{< unsafe >}}
    <script>
        function calculateResources() {
            var seriesInMemory = parseInt(document.getElementById('seriesInMemory').value);
            var samplesPerSecond = parseInt(document.getElementById('samplesPerSecond').value);
            var queriesPerSecond = parseInt(document.getElementById('queriesPerSecond').value);
            var activeSeries = parseInt(document.getElementById('activeSeries').value);
            var firingAlertNotifications = parseInt(document.getElementById('firingAlertNotifications').value);
            var firingAlerts = parseInt(document.getElementById('firingAlerts').value);

            // Distributor
            var distributorCPU = (samplesPerSecond / 25000).toFixed(2);
            var distributorMemory = (samplesPerSecond / 25000).toFixed(2);

            // Ingester
            var ingesterCPU = (seriesInMemory / 300000).toFixed(2);
            var ingesterMemory = (2.5 * seriesInMemory / 300000).toFixed(2);
            var ingesterDisk = (5 * seriesInMemory / 300000).toFixed(2);

            // Query-frontend
            var queryFrontendCPU = (queriesPerSecond / 250).toFixed(2);
            var queryFrontendMemory = (queriesPerSecond / 250).toFixed(2);

            // Query-scheduler
            var querySchedulerCPU = (queriesPerSecond / 500).toFixed(2);
            var querySchedulerMemory = (0.1 * queriesPerSecond / 500).toFixed(2);

            // Querier
            var querierCPU = (queriesPerSecond / 10).toFixed(2);
            var querierMemory = (queriesPerSecond / 10).toFixed(2);

            // Store-gateway
            var storeGatewayCPU = (queriesPerSecond / 10).toFixed(2);
            var storeGatewayMemory = (queriesPerSecond / 10).toFixed(2);
            var storeGatewayDisk = (13 * activeSeries / 1000000).toFixed(2);

            // Ruler
            var rulerCPU = querierCPU; // Same as Querier
            var rulerMemory = querierMemory; // Same as Querier

            // Compactor
            var compactorCPU = (activeSeries/ 20000000).toFixed(2);
            var compactorMemory = (4*activeSeries/ 20000000).toFixed(2);
            var compactorDisk = (300*activeSeries/ 20000000).toFixed(2);

            // Alertmanager
            var alertmanagerCPU = (firingAlertNotifications / 100).toFixed(2);
            var alertmanagerMemory = (firingAlerts / 5000).toFixed(2);

                // Declare a variable for total memory
                var totalMemory = 0;
            // Declare variables for total CPU and Disk usage
                var totalCPU = 0;
                var totalDisk = 0;

            totalCPU += parseFloat(distributorCPU);
                totalCPU += parseFloat(ingesterCPU);
                totalCPU += parseFloat(queryFrontendCPU);
                totalCPU += parseFloat(querySchedulerCPU);
                totalCPU += parseFloat(querierCPU);
                totalCPU += parseFloat(storeGatewayCPU);
                totalCPU += parseFloat(rulerCPU);
                totalCPU += parseFloat(compactorCPU);
                totalCPU += parseFloat(alertmanagerCPU);

                totalDisk += parseFloat(ingesterDisk);
                totalDisk += parseFloat(storeGatewayDisk);
                totalDisk += parseFloat(compactorDisk);

                // Calculate individual components' resources and add to total memory
                totalMemory += parseFloat(distributorMemory);
                totalMemory += parseFloat(ingesterMemory);
                totalMemory += parseFloat(queryFrontendMemory);
                totalMemory += parseFloat(querySchedulerMemory);
                totalMemory += parseFloat(querierMemory);
                totalMemory += parseFloat(storeGatewayMemory);
                totalMemory += parseFloat(rulerMemory);
                totalMemory += parseFloat(compactorMemory);
                totalMemory += parseFloat(alertmanagerMemory);


            // Set calculated values
            document.getElementById('distributorCPU').textContent = distributorCPU;
            document.getElementById('distributorMemory').textContent = distributorMemory;
            document.getElementById('ingesterCPU').textContent = ingesterCPU;
            document.getElementById('ingesterMemory').textContent = ingesterMemory;
            document.getElementById('ingesterDisk').textContent = ingesterDisk;
            document.getElementById('queryFrontendCPU').textContent = queryFrontendCPU;
            document.getElementById('queryFrontendMemory').textContent = queryFrontendMemory;
            document.getElementById('querySchedulerCPU').textContent = querySchedulerCPU;
            document.getElementById('querySchedulerMemory').textContent = querySchedulerMemory;
            document.getElementById('querierCPU').textContent = querierCPU;
            document.getElementById('querierMemory').textContent = querierMemory;
            document.getElementById('storeGatewayCPU').textContent = storeGatewayCPU;
            document.getElementById('storeGatewayMemory').textContent = storeGatewayMemory;
            document.getElementById('storeGatewayDisk').textContent = storeGatewayDisk;
            document.getElementById('rulerCPU').textContent = rulerCPU;
            document.getElementById('rulerMemory').textContent = rulerMemory;
            document.getElementById('compactorCPU').textContent = compactorCPU;
            document.getElementById('compactorMemory').textContent = compactorMemory;
            document.getElementById('compactorDisk').textContent = compactorDisk;
            document.getElementById('alertmanagerCPU').textContent = alertmanagerCPU;
            document.getElementById('alertmanagerMemory').textContent = alertmanagerMemory;
            document.getElementById('totalMemory').textContent = totalMemory.toFixed(2)
            document.getElementById('totalCPU').textContent = totalCPU.toFixed(2);
                document.getElementById('totalDisk').textContent = totalDisk.toFixed(2);


            return false; // Prevent form submission
        }

        function setupEventListeners() {
            var fields = ['seriesInMemory', 'samplesPerSecond', 'queriesPerSecond', 'activeSeries', 'firingAlertNotifications', 'firingAlerts'];
            fields.forEach(function(field) {
                document.getElementById(field).addEventListener('change', calculateResources);
                document.getElementById(field).addEventListener('input', calculateResources);
            });
           calculateResources();
        }

        window.onload = setupEventListeners;
    </script>
    <form>
        <label for="activeSeries">Active series:</label>
        <input type="number" id="activeSeries" name="activeSeries" value="1000000"><br><br>

        <label for="seriesInMemory">Series in memory (Active series*Replication factor):</label>
        <input type="number" id="seriesInMemory" name="seriesInMemory" value="3000000"><br><br>

        <label for="seriesInMemory">Samples per second:</label>
        <input type="number" id="samplesPerSecond" name="samplesPerSecond" value="25000"><br><br>

        <label for="queriesPerSecond">Queries per second:</label>
        <input type="number" id="queriesPerSecond" name="queriesPerSecond" value="10"><br><br>

        <label for="firingAlertNotifications">Firing alert notifications:</label>
        <input type="number" id="firingAlertNotifications" name="firingAlertNotifications" value="100"><br><br>

        <label for="firingAlerts">Firing alerts:</label>
        <input type="number" id="firingAlerts" name="firingAlerts" value="5000"><br><br>
    </form>

    <h2>Total Resource Requirements</h2>
    <div>
        <h3>Total Resources</h3>
            <div>CPU: <span id="totalCPU">N/A</span> cores</div>
                <div>Memory: <span id="totalMemory">N/A</span> GB</div>
                    <div>Disk: <span id="totalDisk">N/A</span> GB</div>
    </div>


    <h2>Calculated Requirements</h2>
    <div>
        <h3>Distributor</h3>
        <div>CPU: <span id="distributorCPU">N/A</span> cores</div>
        <div>Memory: <span id="distributorMemory">N/A</span> GB</div>

        <h3>Ingester</h3>
        <div>CPU: <span id="ingesterCPU">N/A</span> cores</div>
        <div>Memory: <span id="ingesterMemory">N/A</span> GB</div>
        <div>Disk: <span id="ingesterDisk">N/A</span> GB</div>

        <h3>Query-frontend</h3>
        <div>CPU: <span id="queryFrontendCPU">N/A</span> cores</div>
        <div>Memory: <span id="queryFrontendMemory">N/A</span> GB</div>

        <h3>Query-scheduler</h3>
        <div>CPU: <span id="querySchedulerCPU">N/A</span> cores</div>
        <div>Memory: <span id="querySchedulerMemory">N/A</span> GB</div>

        <h3>Querier</h3>
        <div>CPU: <span id="querierCPU">N/A</span> cores</div>
        <div>Memory: <span id="querierMemory">N/A</span> GB</div>

        <h3>Store-gateway</h3>
        <div>CPU: <span id="storeGatewayCPU">N/A</span> cores</div>
        <div>Memory: <span id="storeGatewayMemory">N/A</span> GB</div>
        <div>Disk: <span id="storeGatewayDisk">N/A</span> GB</div>

        <h3>Ruler</h3>
        <div>CPU: <span id="rulerCPU">N/A</span> cores</div>
        <div>Memory: <span id="rulerMemory">N/A</span> GB</div>

        <h3>Compactor</h3>
        <div>CPU: <span id="compactorCPU">N/A</span> cores</div>
        <div>Memory: <span id="compactorMemory">N/A</span> GB</div>
        <div>Disk: <span id="compactorDisk">N/A</span> GB</div>

        <h3>Alertmanager</h3>
        <div>CPU: <span id="alertmanagerCPU">N/A</span> cores</div>
        <div>Memory: <span id="alertmanagerMemory">N/A</span> GB</div>
    </div>
{{< /unsafe >}}

