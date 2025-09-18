# Slurm Viewer TUI – TODO List

## 1. Project Setup

* [ ] Initialize Go module: `go mod init github.com/Open-Knowledge-AI/bubble-slurm`
* [ ] Add dependencies:

  * `github.com/charmbracelet/bubbletea`
  * `github.com/charmbracelet/bubbles`
  * `github.com/charmbracelet/lipgloss`
* [ ] Create basic folder structure: `cmd/`, `internal/`, `pkg/`

---

## 2. Basic TUI Skeleton

* [ ] Implement `main.go` with `tea.Model` structure: `Init()`, `Update()`, `View()`
* [ ] Add placeholder table with static job data
* [ ] Implement basic keybindings: `q` to quit, arrow keys to navigate

---

## 3. Data Fetching Layer

* [ ] Implement `getJobs()`:

  * Primary: `squeue --json`
  * Fallback: parse `squeue` CLI output
* [ ] Create async `tea.Cmd` for loading jobs
* [ ] Handle errors and loading states

---

## 4. Job List View

* [ ] Render jobs in scrollable table (`bubbles/table`)
* [ ] Highlight selected row
* [ ] Add job count/status summary in status bar

---

## 5. Job Detail View

* [ ] Implement `Enter` key to open selected job
* [ ] Fetch details via `scontrol show job <id>`
* [ ] Display node allocation, CPUs, memory, GPUs (if any)

---

## 6. Navigation and Keybindings

* [ ] Arrow keys: navigate job list
* [ ] `Enter`: show job detail
* [ ] `r`: refresh job list
* [ ] `f`: filter jobs (by user, partition, state)
* [ ] `c`: cancel job (with confirmation)
* [ ] `q`: quit
* [ ] `g`: toggle GPU view (if enabled)

---

## 7. Filters

* [ ] Implement text input filter (user/partition/state)
* [ ] Apply filters dynamically to job table

---

## 8. Styling

* [ ] Use Lipgloss for:

  * Job state colors (running = green, pending = yellow, failed = red)
  * Highlighting selected rows
  * Status bar styling
* [ ] Optional: add progress bars for memory/CPU/GPU usage

---

## 9. GPU Monitoring (Optional)

* [ ] Detect presence of `nvtop` or `nvidia-smi` at startup
* [ ] Define `GpuUsage` struct:

  ```go
  type GpuUsage struct {
      Node       string
      GpuIndex   int
      MemUsedMB  int
      MemTotalMB int
      Utilization int
      JobID      string
  }
  ```
* [ ] Implement `getGpuMetrics()` to fetch per-GPU memory & utilization
* [ ] Cluster GPU view:

  * Table: Node | GPU ID | Mem Used / Total | Utilization % | Job ID
  * Progress bars for memory usage
* [ ] Job detail GPU view:

  * Show GPUs allocated to the job
  * Memory/Utilization bars per GPU
* [ ] Background refresh for GPU metrics using Bubbletea tick messages
* [ ] Keybindings:

  * `g` to toggle GPU view
  * `←/→` to switch between job list and GPU view
* [ ] Graceful fallback if `nvtop` or `nvidia-smi` not installed

---

## 10. Refresh & Background Updates

* [ ] Implement periodic refresh for jobs and GPU metrics
* [ ] Ensure non-blocking updates (use `tea.Cmd`)

---

## 11. Testing & Deployment

* [ ] Test on HPC login node
* [ ] Verify performance with large job arrays
* [ ] Handle missing commands gracefully (`squeue`, `scontrol`, `nvtop`)

---

## 12. Stretch Goals

* [ ] Node view (`sinfo`) with status, CPUs, memory, GPUs
* [ ] Job history view (`sacct`)
* [ ] Job submission wrapper (`sbatch`) for TUI
* [ ] Config file support (default user, partition, refresh interval)
* [ ] Mouse support (scrolling/selecting rows)
* [ ] GPU filters/sorting by memory usage
* [ ] Sparkline charts for historical GPU usage
