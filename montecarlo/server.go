package montecarlo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/mshindle/structures/ringbuffer"
)

type Visualizer struct {
	calculator *PI
	mux        *http.ServeMux
	server     *http.Server
}

func NewVisualizer(piCalc *PI) *Visualizer {
	v := &Visualizer{
		calculator: piCalc,
		mux:        http.NewServeMux(),
	}
	v.routes()
	return v
}

func (v *Visualizer) routes() {
	v.mux.HandleFunc("/api/data", v.handleDataPoints)
	v.mux.HandleFunc("/dashboard", v.handleDashboard)
}
func (v *Visualizer) handleDataPoints(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Open CORS rules so external local plotting files can connect easily
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Safely pull a point-in-time slice copy out of the rolling buffer.
	// We extract up to the current length of the RingBuffer.
	length := v.calculator.SampleBuffer.Len()
	pointsSnapshot := make([]RenderPoint, 0, length)

	// Pop items or read them into our snapshot array safely
	for i := 0; i < length; i++ {
		var val RenderPoint
		var err error
		if val, err = v.calculator.SampleBuffer.Pop(); err != nil {
			if !errors.Is(err, ringbuffer.ErrEmpty) {
				log.WithError(err).Error("failed to pop sample from buffer")
			}
			break
		}
		pointsSnapshot = append(pointsSnapshot, val)
	}

	total := v.calculator.TotalPoints.Load()
	in := v.calculator.InCircle.Load()
	estimatedPi := 0.0
	if total > 0 {
		estimatedPi = 4.0 * (float64(in) / float64(total))
	}

	// Wrapper payload containing the calculated convergence status
	payload := map[string]any{
		"estimated_pi": estimatedPi,
		"total_points": total,
		"samples":      pointsSnapshot,
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to serialize visualization data", http.StatusInternalServerError)
	}
}

// internal/montecarlo/server.go

// handleDashboard returns a zero-dependency HTML dashboard using Apache ECharts via CDN
func (v *Visualizer) handleDashboard(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Monte Carlo PI Visualization</title>
		<script src="https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/echarts.min.js"></script>
		<style>
			body { font-family: sans-serif; background: #f9f9f9; text-align: center; }
			#chart { width: 600px; height: 600px; margin: 20px auto; background: white; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
			#stats { font-size: 1.2rem; color: #333; margin-top: 20px; }
		</style>
	</head>
	<body>
		<h1>Monte Carlo Simulation Engine</h1>
		<div id="stats">Estimating PI...</div>
		<div id="chart"></div>

		<script>
			var chartDom = document.getElementById('chart');
			var myChart = echarts.init(chartDom);
			
			// Store the points globally in the browser
			var inCircleData = [];
			var outCircleData = [];
			var MAX_VISIBLE_POINTS = 2500; // Protect the browser memory
			
			function updateDashboard() {
				fetch('/api/data')
				.then(res => res.json())
				.then(data => {
					document.getElementById('stats').innerText = 
							"Total Computed: " + data.total_points + " | Estimated PI: " + data.estimated_pi.toFixed(6);
			
					// Defensive check in case the buffer was empty on this exact tick
					var samples = data.samples || [];

					// Accumulate the newly drained points
					samples.forEach(p => {
						if (p.in_circle) {
							inCircleData.push([p.x, p.y]);
						} else {
							outCircleData.push([p.x, p.y]);
						}
					});
	
					// Shift older points off the front to prevent the browser from crashing
					while ((inCircleData.length + outCircleData.length) > MAX_VISIBLE_POINTS) {
						if (inCircleData.length > outCircleData.length) {
							inCircleData.shift();
						} else {
							outCircleData.shift();
						}
					}
	
					myChart.setOption({
						// EXPLICITLY set type to 'value' so Echarts knows how to plot the X/Y floats
						xAxis: { type: 'value', min: -1, max: 1, splitLine: { show: false } },
						yAxis: { type: 'value', min: -1, max: 1, splitLine: { show: false } },
						legend: { data: ['In Circle', 'Outside Circle'], bottom: 10 },
						animation: false, // Turn off animation for high-performance rendering
						series: [
							{ type: 'scatter', name: 'In Circle', data: inCircleData, itemStyle: { color: '#ef4444' }, symbolSize: 4 },
							{ type: 'scatter', name: 'Outside Circle', data: outCircleData, itemStyle: { color: '#3b82f6' }, symbolSize: 4 }
						]
					});
				}).catch(err => console.error("Error fetching telemetry:", err));
			}

			setInterval(updateDashboard, 500);
			updateDashboard();
		</script>
	</body>
	</html>`
	_, _ = w.Write([]byte(html))
}

func (v *Visualizer) ListenAndServe(addr string) error {
	log.Infof("Visualization web interface starting on http://localhost%s/dashboard", addr)
	v.server = &http.Server{
		Addr:         addr,
		Handler:      v.mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return v.server.ListenAndServe()
}

func (v *Visualizer) Shutdown(ctx context.Context) error {
	if v.server == nil {
		return nil
	}
	return v.server.Shutdown(ctx)
}
