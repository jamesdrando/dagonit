package graph

import (
	"fmt"
	"net/http"
	"time"
)

// StartServer starts a web server to visualize the graph.
func StartServer(port int, g *Graph) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(indexHTML))
	})

	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		data, err := g.ExportToJSON()
		if err != nil {
			http.Error(w, "Failed to export graph", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Serving graph visualization at http://localhost%s\n", addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.ListenAndServe()
}

const indexHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dagonit Graph Visualizer</title>
    <style>
        body {
            margin: 0;
            overflow: hidden;
            background-color: #121212; /* Deep matte black */
            color: #e0e0e0;
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        }
        #canvas {
            display: block;
            width: 100vw;
            height: 100vh;
            cursor: grab;
        }
        #canvas:active {
            cursor: grabbing;
        }
        #controls {
            position: absolute;
            top: 20px;
            left: 20px;
            background: rgba(30, 30, 30, 0.8);
            padding: 15px;
            border-radius: 8px;
            backdrop-filter: blur(5px);
            border: 1px solid #333;
            box-shadow: 0 4px 6px rgba(0,0,0,0.3);
            display: flex;
            flex-direction: column;
            gap: 10px;
            min-width: 200px;
        }
        .control-group {
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        button {
            background: #2d2d2d;
            border: 1px solid #444;
            color: #ccc;
            padding: 5px 10px;
            border-radius: 4px;
            cursor: pointer;
            transition: background 0.2s;
        }
        button:hover {
            background: #3d3d3d;
            color: #fff;
        }
        label {
            font-size: 12px;
            color: #aaa;
        }
        #status {
            font-size: 11px;
            color: #666;
            margin-top: 5px;
        }
        #tooltip {
            position: absolute;
            pointer-events: none;
            background: rgba(0, 0, 0, 0.9);
            border: 1px solid #444;
            padding: 8px 12px;
            border-radius: 4px;
            font-size: 12px;
            display: none;
            z-index: 10;
            max-width: 300px;
            color: #eee;
        }
        .node-type {
            font-weight: bold;
            color: #4a90e2;
            margin-bottom: 2px;
            display: block;
            text-transform: uppercase;
            font-size: 10px;
        }
    </style>
</head>
<body>
        <div id="controls">
            <div class="control-group">
                <span style="font-weight:bold; color:#fff;">Dagonit</span>
                <span id="node-count" style="font-size:11px; color:#888;">0 nodes</span>
            </div>
            <div class="control-group">
                <button id="btn-reset">Reset View</button>
                <button id="btn-pause">Pause</button>
            </div>
            <div class="control-group">
                <label>View Mode:</label>
                <select id="dimension-select">
                    <option value="full">Full Graph</option>
                    <option value="files">Files Only</option>
                    <option value="symbols">Symbols Only</option>
                    <option value="modules">Modules Only</option>
                    <option value="dependencies">Dependencies</option>
                    <option value="calls">Call Graphs</option>
                    <option value="scc">SCCs</option>
                    <option value="custom">Custom</option>
                </select>
            </div>
            <div class="control-group" id="custom-controls" style="display: none;">
                <div style="font-size: 11px; color: #aaa; margin-bottom: 5px;">Node Types:</div>
                <div>
                    <label><input type="checkbox" id="node-file" checked> Files</label>
                    <label><input type="checkbox" id="node-symbol" checked> Symbols</label>
                    <label><input type="checkbox" id="node-module" checked> Modules</label>
                    <label><input type="checkbox" id="node-scc" checked> SCCs</label>
                    <label><input type="checkbox" id="node-repository"> Repository</label>
                    <label><input type="checkbox" id="node-test"> Tests</label>
                    <label><input type="checkbox" id="node-build-target"> Build Targets</label>
                    <label><input type="checkbox" id="node-config"> Configs</label>
                </div>
                <div style="font-size: 11px; color: #aaa; margin: 10px 0 5px 0;">Edge Types:</div>
                <div>
                    <label><input type="checkbox" id="edge-contains" checked> Contains</label>
                    <label><input type="checkbox" id="edge-imports" checked> Imports</label>
                    <label><input type="checkbox" id="edge-calls" checked> Calls</label>
                    <label><input type="checkbox" id="edge-tested-by" checked> Tested By</label>
                    <label><input type="checkbox" id="edge-builds" checked> Builds</label>
                    <label><input type="checkbox" id="edge-configures" checked> Configures</label>
                </div>
            </div>
            <div class="control-group" id="legend" style="font-size: 10px; color: #aaa;">
                <strong>Legend:</strong><br>
                <span style="color: #4a90e2;">■</span> File &nbsp;
                <span style="color: #e24a4a;">■</span> Module &nbsp;
                <span style="color: #f5a623;">■</span> Package/SCC &nbsp;
                <span style="color: #7bdcb5;">■</span> Symbol &nbsp;
                <span style="color: #9d4edd;">■</span> Repository &nbsp;
                <span style="color: #f8e71c;">■</span> Test &nbsp;
                <span style="color: #ff9ff3;">■</span> Build Target &nbsp;
                <br>
                <span style="color: rgba(74, 144, 226, 0.4);">―</span> Contains &nbsp;
                <span style="color: rgba(226, 74, 74, 0.4);">―</span> Imports &nbsp;
                <span style="color: rgba(123, 220, 181, 0.4);">―</span> Calls &nbsp;
                <span style="color: rgba(245, 166, 35, 0.4);">―</span> Tested By &nbsp;
                <span style="color: rgba(248, 231, 28, 0.4);">―</span> Builds &nbsp;
                <span style="color: rgba(255, 159, 243, 0.4);">―</span> Configures
            </div>
            <div class="control-group">
                <label>Repulsion</label>
                <input type="range" id="param-repulsion" min="10" max="1000" value="200">
            </div>
            <div class="control-group">
                <label>Spring</label>
                <input type="range" id="param-spring" min="0.001" max="0.1" step="0.001" value="0.01">
            </div>
            <div id="status">Simulating...</div>
        </div>
    <div id="tooltip"></div>
    <canvas id="canvas"></canvas>

    <script>
        const CONFIG = {
            colors: {
                background: '#121212',
                node: {
                    default: '#3d3d3d',
                    file: '#4a90e2',      
                    module: '#e24a4a',    
                    package: '#f5a623',   
                    symbol: '#7bdcb5',    
                    repository: '#9d4edd',
                    test: '#f8e71c',
                    build_target: '#ff9ff3',
                    scc: '#f5a623',
                    highlight: '#ffffff',
                    text: '#cccccc'
                },
                edge: {
                    contains: 'rgba(74, 144, 226, 0.4)',
                    imports: 'rgba(226, 74, 74, 0.4)',
                    calls: 'rgba(123, 220, 181, 0.4)',
                    tested_by: 'rgba(245, 166, 35, 0.4)',
                    builds: 'rgba(248, 231, 28, 0.4)',
                    configures: 'rgba(255, 159, 243, 0.4)',
                    default: 'rgba(80, 80, 80, 0.4)',
                    highlight: 'rgba(200, 200, 200, 0.8)'
                }
            },
            physics: {
                repulsion: 200,
                springLength: 100,
                springStrength: 0.02,
                damping: 0.85, 
                centerGravity: 0.0005,
                maxVelocity: 5
            }
        };

        let nodes = [];
        let edges = [];
        let nodeMap = {}; 
        let adj = {};     

        let camera = { x: 0, y: 0, zoom: 1 };
        let isDragging = false;
        let isPanning = false;
        let dragStart = { x: 0, y: 0 };
        let draggedNode = null;
        let hoveredNode = null;
        let animating = true;
        
        const canvas = document.getElementById('canvas');
        const ctx = canvas.getContext('2d');
        const tooltip = document.getElementById('tooltip');
        
        function resize() {
            canvas.width = window.innerWidth;
            canvas.height = window.innerHeight;
            if (camera.x === 0 && camera.y === 0) {
                camera.x = canvas.width / 2;
                camera.y = canvas.height / 2;
            }
        }
        window.addEventListener('resize', resize);
        resize();

        fetch('/data')
            .then(res => res.json())
            .then(data => {
                initGraph(data);
                applyFilters(); // Apply initial filters
                requestAnimationFrame(loop);
            })
            .catch(err => {
                console.error("Failed to load graph:", err);
                document.getElementById('status').innerText = "Error loading data.";
            });

        function initGraph(data) {
            nodes = data.nodes.map(n => ({
                id: n.id,
                type: n.type,
                metadata: n.metadata,
                x: (Math.random() - 0.5) * 800, 
                y: (Math.random() - 0.5) * 800,
                vx: 0,
                vy: 0,
                radius: getNodeRadius(n.type)
            }));

            nodes.forEach(n => { nodeMap[n.id] = n; adj[n.id] = []; });

            edges = data.edges.map(e => {
                const source = nodeMap[e.from];
                const target = nodeMap[e.to];
                if (source && target) {
                    const edge = { source, target, type: e.type };
                    adj[source.id].push(edge);
                    adj[target.id].push(edge); 
                    return edge;
                }
                return null;
            }).filter(e => e !== null);

            document.getElementById('node-count').innerText = nodes.length + " nodes";
            
            const angleStep = 0.5;
            nodes.forEach((n, i) => {
                const angle = i * angleStep;
                const r = 50 + (i * 10); 
                n.x = Math.cos(angle) * r;
                n.y = Math.sin(angle) * r;
            });
        }

        function getNodeRadius(type) {
            switch(type) {
                case 'module': return 12;
                case 'package': return 10;
                case 'file': return 8; 
                case 'symbol': return 6;
                case 'repository': return 14;
                case 'test': return 7;
                case 'build_target': return 9;
                case 'scc': return 11;
                default: return 5;
            }
        }

        function getNodeColor(type) {
            switch(type) {
                case 'module': return CONFIG.colors.node.module;
                case 'package': return CONFIG.colors.node.package;
                case 'file': return CONFIG.colors.node.file;
                case 'symbol': return CONFIG.colors.node.symbol;
                case 'repository': return CONFIG.colors.node.repository;
                case 'test': return CONFIG.colors.node.test;
                case 'build_target': return CONFIG.colors.node.build_target;
                case 'scc': return CONFIG.colors.node.scc;
                default: return CONFIG.colors.node.default;
            }
        }

        function updatePhysics() {
            if (!animating) return;

            for (let i = 0; i < nodes.length; i++) {
                let n1 = nodes[i];
                for (let j = i + 1; j < nodes.length; j++) {
                    let n2 = nodes[j];
                    let dx = n1.x - n2.x;
                    let dy = n1.y - n2.y;
                    let distSq = dx * dx + dy * dy;
                    if (distSq === 0) distSq = 0.1;
                    
                    if (distSq > 500000) continue; 

                    let dist = Math.sqrt(distSq);
                    let force = CONFIG.physics.repulsion * 50 / distSq; 
                    
                    let fx = (dx / dist) * force;
                    let fy = (dy / dist) * force;

                    n1.vx += fx;
                    n1.vy += fy;
                    n2.vx -= fx;
                    n2.vy -= fy;
                }
            }

            for (let e of edges) {
                let n1 = e.source;
                let n2 = e.target;
                let dx = n2.x - n1.x;
                let dy = n2.y - n1.y;
                let dist = Math.sqrt(dx * dx + dy * dy);
                if (dist === 0) dist = 0.1;

                let displacement = dist - CONFIG.physics.springLength;
                let force = displacement * CONFIG.physics.springStrength;

                let fx = (dx / dist) * force;
                let fy = (dy / dist) * force;

                n1.vx += fx;
                n1.vy += fy;
                n2.vx -= fx;
                n2.vy -= fy;
            }

            for (let n of nodes) {
                if (n === draggedNode) continue;

                n.vx -= n.x * CONFIG.physics.centerGravity;
                n.vy -= n.y * CONFIG.physics.centerGravity;

                n.vx *= CONFIG.physics.damping;
                n.vy *= CONFIG.physics.damping;

                let speed = Math.sqrt(n.vx * n.vx + n.vy * n.vy);
                if (speed > CONFIG.physics.maxVelocity) {
                    n.vx = (n.vx / speed) * CONFIG.physics.maxVelocity;
                    n.vy = (n.vy / speed) * CONFIG.physics.maxVelocity;
                }

                n.x += n.vx;
                n.y += n.vy;
            }
        }

        function draw() {
            ctx.fillStyle = CONFIG.colors.background;
            ctx.fillRect(0, 0, canvas.width, canvas.height);

            ctx.save();
            ctx.translate(camera.x, camera.y);
            ctx.scale(camera.zoom, camera.zoom);

            ctx.lineWidth = 1;
            for (let e of edges) {
                let isHovered = hoveredNode && (e.source === hoveredNode || e.target === hoveredNode);
                ctx.strokeStyle = isHovered ? CONFIG.colors.edge.highlight : CONFIG.colors.edge.default;
                ctx.globalAlpha = isHovered ? 0.8 : 0.2;
                
                ctx.beginPath();
                ctx.moveTo(e.source.x, e.source.y);
                ctx.lineTo(e.target.x, e.target.y);
                ctx.stroke();
            }
            ctx.globalAlpha = 1.0;

            for (let n of nodes) {
                let isHovered = n === hoveredNode;
                
                ctx.beginPath();
                ctx.arc(n.x, n.y, n.radius, 0, Math.PI * 2);
                ctx.fillStyle = isHovered ? CONFIG.colors.node.highlight : getNodeColor(n.type);
                ctx.fill();

                if (isHovered) {
                    ctx.strokeStyle = '#fff';
                    ctx.lineWidth = 2;
                    ctx.stroke();
                }

                if (camera.zoom > 0.6 || isHovered || n.type === 'module') {
                    ctx.fillStyle = isHovered ? '#fff' : CONFIG.colors.node.text;
                    ctx.font = '10px sans-serif';
                    ctx.textAlign = 'center';
                    
                    let label = n.id.split('/').pop();
                    if (label.length > 20 && !isHovered) label = label.substr(0, 17) + "...";
                    ctx.fillText(label, n.x, n.y + n.radius + 12);
                }
            }

            ctx.restore();
        }

        function loop() {
            updatePhysics();
            draw();
            requestAnimationFrame(loop);
        }

        function toWorld(x, y) {
            return {
                x: (x - camera.x) / camera.zoom,
                y: (y - camera.y) / camera.zoom
            };
        }

        function getHoveredNode(mx, my) {
            let worldPos = toWorld(mx, my);
            for (let i = nodes.length - 1; i >= 0; i--) {
                let n = nodes[i];
                let dx = worldPos.x - n.x;
                let dy = worldPos.y - n.y;
                if (dx*dx + dy*dy < (n.radius + 5) * (n.radius + 5)) {
                    return n;
                }
            }
            return null;
        }

        canvas.addEventListener('mousedown', e => {
            const node = getHoveredNode(e.clientX, e.clientY);
            if (e.button === 0) { 
                if (node) {
                    draggedNode = node;
                    draggedNode.vx = 0;
                    draggedNode.vy = 0;
                    isDragging = true; 
                } else {
                    isPanning = true; 
                    dragStart.x = e.clientX - camera.x;
                    dragStart.y = e.clientY - camera.y;
                    canvas.style.cursor = 'grabbing';
                }
            }
        });

        window.addEventListener('mousemove', e => {
            const node = getHoveredNode(e.clientX, e.clientY);
            
            if (node !== hoveredNode) {
                hoveredNode = node;
                if (node) {
                    tooltip.style.display = 'block';
                    tooltip.innerHTML = "<span class='node-type'>" + node.type + "</span>" + node.id;
                    canvas.style.cursor = 'pointer';
                } else {
                    tooltip.style.display = 'none';
                    canvas.style.cursor = isPanning ? 'grabbing' : 'default';
                }
            }
            
            if (hoveredNode) {
                tooltip.style.left = (e.clientX + 15) + 'px';
                tooltip.style.top = (e.clientY + 15) + 'px';
            }

            if (draggedNode) {
                let pos = toWorld(e.clientX, e.clientY);
                draggedNode.x = pos.x;
                draggedNode.y = pos.y;
            } else if (isPanning) {
                camera.x = e.clientX - dragStart.x;
                camera.y = e.clientY - dragStart.y;
            }
        });

        window.addEventListener('mouseup', () => {
            draggedNode = null;
            isPanning = false;
            isDragging = false;
            canvas.style.cursor = 'default';
        });

        canvas.addEventListener('wheel', e => {
            e.preventDefault();
            const zoomSensitivity = 0.001;
            const delta = -e.deltaY * zoomSensitivity;
            const oldZoom = camera.zoom;
            let newZoom = Math.max(0.1, Math.min(5, oldZoom + delta));
            
            const mouseX = e.clientX;
            const mouseY = e.clientY;
            
            const worldX = (mouseX - camera.x) / oldZoom;
            const worldY = (mouseY - camera.y) / oldZoom;
            
            camera.x = mouseX - worldX * newZoom;
            camera.y = mouseY - worldY * newZoom;
            camera.zoom = newZoom;
        });

        document.getElementById('btn-reset').addEventListener('click', () => {
             camera.x = canvas.width / 2;
             camera.y = canvas.height / 2;
             camera.zoom = 1;
        });

        document.getElementById('btn-pause').addEventListener('click', (e) => {
            animating = !animating;
            e.target.innerText = animating ? "Resume" : "Pause";
        });

        document.getElementById('param-repulsion').addEventListener('input', (e) => {
            CONFIG.physics.repulsion = parseFloat(e.target.value);
        });
        document.getElementById('param-spring').addEventListener('input', (e) => {
            CONFIG.physics.springStrength = parseFloat(e.target.value);
        });

        // Dimension filtering controls
        const dimensionSelect = document.getElementById('dimension-select');
        const customControls = document.getElementById('custom-controls');
        
        dimensionSelect.addEventListener('change', (e) => {
            const mode = e.target.value;
            // Show custom controls only when custom mode is selected
            customControls.style.display = mode === 'custom' ? 'block' : 'none';
            applyPresetFilters(mode);
            applyFilters(); // Reapply filters with new preset
        });
        
        // Apply preset filters when checkboxes change in custom mode
        document.querySelectorAll('#custom-controls input[type="checkbox"]').forEach(checkbox => {
            checkbox.addEventListener('change', applyFilters);
        });

        // Filtering state
        let filterState = {
            nodeTypes: new Set(['file', 'symbol', 'module', 'scc', 'repository', 'test', 'build_target', 'config']),
            edgeTypes: new Set(['contains', 'imports', 'calls', 'tested_by', 'builds', 'configures']),
            activePreset: 'full'
        };

        // Apply preset filters
        function applyPresetFilters(mode) {
            // Reset to all types
            filterState.nodeTypes = new Set(['file', 'symbol', 'module', 'scc', 'repository', 'test', 'build_target', 'config']);
            filterState.edgeTypes = new Set(['contains', 'imports', 'calls', 'tested_by', 'builds', 'configures']);
            
            switch(mode) {
                case 'files':
                    filterState.nodeTypes = new Set(['file']);
                    filterState.edgeTypes = new Set(['contains']);
                    break;
                case 'symbols':
                    filterState.nodeTypes = new Set(['symbol']);
                    filterState.edgeTypes = new Set(['calls']);
                    break;
                case 'modules':
                    filterState.nodeTypes = new Set(['module']);
                    filterState.edgeTypes = new Set(['contains']);
                    break;
                case 'dependencies':
                    filterState.nodeTypes = new Set(['file', 'module']);
                    filterState.edgeTypes = new Set(['imports', 'contains']);
                    break;
                case 'calls':
                    filterState.nodeTypes = new Set(['symbol']);
                    filterState.edgeTypes = new Set(['calls']);
                    break;
                case 'scc':
                    filterState.nodeTypes = new Set(['scc']);
                    filterState.edgeTypes = new Set([]); // SCCs typically don't have outgoing edges in condensed graph
                    break;
                case 'full':
                case 'custom':
                    // Keep all types (will be refined by checkboxes in custom mode)
                    break;
            }
            
            filterState.activePreset = mode;
            
            // Update checkboxes to match preset (except in custom mode where we preserve manual selections)
            if (mode !== 'custom') {
                updateCheckboxesFromFilterState();
            }
        }

        // Update checkboxes to match current filter state
        function updateCheckboxesFromFilterState() {
            document.getElementById('node-file').checked = filterState.nodeTypes.has('file');
            document.getElementById('node-symbol').checked = filterState.nodeTypes.has('symbol');
            document.getElementById('node-module').checked = filterState.nodeTypes.has('module');
            document.getElementById('node-scc').checked = filterState.nodeTypes.has('scc');
            document.getElementById('node-repository').checked = filterState.nodeTypes.has('repository');
            document.getElementById('node-test').checked = filterState.nodeTypes.has('test');
            document.getElementById('node-build-target').checked = filterState.nodeTypes.has('build_target');
            document.getElementById('node-config').checked = filterState.nodeTypes.has('config');
            
            document.getElementById('edge-contains').checked = filterState.edgeTypes.has('contains');
            document.getElementById('edge-imports').checked = filterState.edgeTypes.has('imports');
            document.getElementById('edge-calls').checked = filterState.edgeTypes.has('calls');
            document.getElementById('edge-tested-by').checked = filterState.edgeTypes.has('tested_by');
            document.getElementById('edge-builds').checked = filterState.edgeTypes.has('builds');
            document.getElementById('edge-configures').checked = filterState.edgeTypes.has('configures');
        }

        // Apply filters to nodes and edges
        function applyFilters() {
            // Update node visibility
            nodes.forEach(node => {
                const typeMatch = filterState.nodeTypes.has(node.type);
                node.visible = typeMatch;
            });
            
            // Update edge visibility based on source and target node visibility and edge type
            edges.forEach(edge => {
                const typeMatch = filterState.edgeTypes.has(edge.type);
                const sourceVisible = edge.source.visible;
                const targetVisible = edge.target.visible;
                edge.visible = typeMatch && sourceVisible && targetVisible;
            });
            
            // Update node count display
            const visibleNodes = nodes.filter(node => node.visible).length;
            document.getElementById('node-count').innerText = visibleNodes + " nodes";
        }

        // Modify draw function to respect visibility
        const originalDraw = draw;
        draw = function() {
            // Call original draw but skip invisible nodes/edges
            ctx.fillStyle = CONFIG.colors.background;
            ctx.fillRect(0, 0, canvas.width, canvas.height);

            ctx.save();
            ctx.translate(camera.x, camera.y);
            ctx.scale(camera.zoom, camera.zoom);

            ctx.lineWidth = 1;
            // Draw only visible edges
            for (let e of edges) {
                if (!e.visible) continue;
                
                let isHovered = hoveredNode && (e.source === hoveredNode || e.target === hoveredNode);
                ctx.strokeStyle = isHovered ? CONFIG.colors.edge.highlight : CONFIG.colors.edge.default;
                ctx.globalAlpha = isHovered ? 0.8 : 0.2;
                
                ctx.beginPath();
                ctx.moveTo(e.source.x, e.source.y);
                ctx.lineTo(e.target.x, e.target.y);
                ctx.stroke();
            }
            ctx.globalAlpha = 1.0;

            // Draw only visible nodes
            for (let n of nodes) {
                if (!n.visible) continue;
                
                let isHovered = n === hoveredNode;
                
                ctx.beginPath();
                ctx.arc(n.x, n.y, n.radius, 0, Math.PI * 2);
                ctx.fillStyle = isHovered ? CONFIG.colors.node.highlight : getNodeColor(n.type);
                ctx.fill();

                if (isHovered) {
                    ctx.strokeStyle = '#fff';
                    ctx.lineWidth = 2;
                    ctx.stroke();
                }

                if (camera.zoom > 0.6 || isHovered || n.type === 'module') {
                    ctx.fillStyle = isHovered ? '#fff' : CONFIG.colors.node.text;
                    ctx.font = '10px sans-serif';
                    ctx.textAlign = 'center';
                    
                    let label = n.id.split('/').pop();
                    if (label.length > 20 && !isHovered) label = label.substr(0, 17) + "...";
                    ctx.fillText(label, n.x, n.y + n.radius + 12);
                }
            }

            ctx.restore();
        };
    </script>
</body>
</html>
`
