let TOKEN = localStorage.getItem("token") || "";
        const GATEWAY = "http://localhost:8000/api/v1";

        // Auto-check session on load
        if (TOKEN) showApp();

        async function login() {
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            const btn = document.getElementById('login-btn');

            btn.innerHTML = "Authenticating...";
            btn.disabled = true;

            try {
                const res = await fetch(`${GATEWAY}/identity/login`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email, password })
                });

                if (res.ok) {
                    const data = await res.json();
                    TOKEN = data.token;
                    localStorage.setItem("token", TOKEN); // Save session
                    showApp();
                } else {
                    alert("Login failed! Check credentials.");
                }
            } catch (err) {
                alert("Gateway connection error.");
            } finally {
                btn.innerHTML = "Sign In";
                btn.disabled = false;
            }
        }

        function showApp() {
            document.getElementById('login-section').classList.add('hidden');
            document.getElementById('app-section').classList.remove('hidden');
            document.getElementById('user-info').classList.remove('hidden');
            search();
            loadOrders();
        }

        function logout() {
            localStorage.clear();
            location.reload();
        }

        async function search() {
            const text = document.getElementById('search-query').value;
            // Use the "Smart Search" pattern: send q to the backend
            const res = await fetch(`${GATEWAY}/search?q=${text}`);
            const data = await res.json();

            const html = data.items.map(p => `
                <div class="bg-white p-5 rounded-2xl border border-slate-200 hover:border-orange-300 shadow-sm transition group">
                    <div class="flex justify-between items-start mb-4">
                        <div>
                            <span class="text-xs font-bold text-orange-600 uppercase tracking-tighter">${p.category || 'General'}</span>
                            <h4 class="text-lg font-bold text-slate-800">${p.name}</h4>
                        </div>
                        <span class="text-xl font-bold">$${p.price}</span>
                    </div>
                    <p class="text-sm text-slate-500 mb-6 line-clamp-2">${p.description || 'No description available.'}</p>
                    <button onclick="placeOrder('${p.id}')" class="w-full py-2 bg-slate-50 text-slate-900 font-bold rounded-lg group-hover:bg-orange-600 group-hover:text-white transition">
                        Order Now
                    </button>
                </div>
            `).join('');
            document.getElementById('product-list').innerHTML = html || "<p class='text-slate-400'>No products found matching your search.</p>";
        }

        async function placeOrder(productId) {
            const res = await fetch(`${GATEWAY}/orders`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${TOKEN}`
                },
                body: JSON.stringify({ productId, quantity: 1 })
            });
            if (res.ok) {
                loadOrders();
            } else if (res.status === 401) {
                logout();
            }
        }

        async function loadOrders() {
            const res = await fetch(`${GATEWAY}/orders`, {
                headers: { 'Authorization': `Bearer ${TOKEN}` }
            });
            if (!res.ok) return;
            const data = await res.json();

            document.getElementById('order-count').innerText = data.length;
            document.getElementById('order-history').innerHTML = data.map(o => `
            <div class="bg-slate-800 p-4 rounded-xl border border-slate-700 relative group">
            <!-- THE CROSS BUTTON -->
            ${o.status !== 'PENDING' ?
                    `<button onclick="updateStatus('${o.id}', 'ARCHIVED')" 
                    class="absolute -top-2 -right-2 bg-red-500 text-white w-6 h-6 rounded-full text-xs opacity-0 group-hover:opacity-100 transition">
                    ✕
                </button>` : ''}

            <div class="flex justify-between text-xs text-slate-400 mb-2">
                <span>#${o.id.substring(0, 8)}</span>
                <span class="${o.status === 'COMPLETED' ? 'text-green-400' : 'text-orange-400'} font-bold">${o.status}</span>
            </div>
            
            <div class="flex justify-between items-center">
                <span class="font-bold text-lg">$${o.totalPrice}</span>
                ${o.status === 'PENDING' ?
                    `<button onclick="updateStatus('${o.id}', 'COMPLETED')" class="text-xs bg-white text-slate-900 px-3 py-1 rounded-full font-bold">Mark Delivered</button>`
                    : ''}
            </div>
        </div>
    `).join('');
        }

        async function updateStatus(orderId, newStatus) {
            await fetch(`${GATEWAY}/orders/${orderId}/status`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${TOKEN}`
                },
                body: JSON.stringify({ status: newStatus })
            });
            loadOrders();
        }