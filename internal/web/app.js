let TOKEN = localStorage.getItem("token") || "";
const GATEWAY = "http://localhost:8000/api/v1";

// On Page Load: Check if we have a session
if (TOKEN) {
    showApp();
}

/**
 * Switch between Login and Register forms
 */
function toggleAuth(showRegister) {
    document.getElementById('login-form').classList.toggle('hidden', showRegister);
    document.getElementById('register-form').classList.toggle('hidden', !showRegister);
}

/**
 * Handle Registration
 */
async function register() {
    const username = document.getElementById('reg-username').value;
    const email = document.getElementById('reg-email').value;
    const password = document.getElementById('reg-password').value;
    const btn = document.getElementById('register-btn');

    if (!username || !email || !password) return alert("Please fill all fields");

    btn.innerHTML = "Creating account...";
    try {
        const res = await fetch(`${GATEWAY}/identity/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, email, password })
        });
        if (res.ok) {
            alert("Success! Now please sign in.");
            toggleAuth(false);
        } else {
            alert("Registration failed. Email might exist.");
        }
    } catch (err) { alert("Gateway error."); }
    btn.innerHTML = "Create Account";
}

/**
 * Handle Login
 */
async function login() {
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const btn = document.getElementById('login-btn');

    btn.innerHTML = "Authenticating...";
    try {
        const res = await fetch(`${GATEWAY}/identity/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });
        if (res.ok) {
            const data = await res.json();
            TOKEN = data.token;
            localStorage.setItem("token", TOKEN);
            showApp();
        } else {
            alert("Invalid credentials.");
        }
    } catch (err) { alert("Check if Gateway is running on port 8000"); }
    btn.innerHTML = "Sign In";
}

/**
 * Setup the UI for Logged-In state
 */
function showApp() {
    document.getElementById('auth-section').classList.add('hidden');
    document.getElementById('app-section').classList.remove('hidden');
    document.getElementById('user-info').classList.remove('hidden');
    search();
    loadOrders();
}

function logout() {
    localStorage.clear();
    location.reload();
}

/**
 * Search logic using the "Smart Search" Pattern
 */
async function search() {
    const q = document.getElementById('search-query').value;
    try {
        const res = await fetch(`${GATEWAY}/search?q=${q}`);
        const data = await res.json();
        
        // Map over data.items because we added Pagination in Go
        const html = (data.items || []).map(p => `
            <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm order-card">
                <div class="flex justify-between items-start mb-4">
                    <div>
                        <span class="text-[10px] font-bold text-orange-600 uppercase tracking-widest">${p.category || 'General'}</span>
                        <h4 class="text-lg font-bold text-slate-800">${p.name}</h4>
                    </div>
                    <span class="text-xl font-bold text-slate-900">$${p.price}</span>
                </div>
                <button onclick="placeOrder('${p.id}')" class="w-full py-2 bg-orange-600 text-white font-bold rounded-lg hover:bg-orange-700 transition">
                    Order Now
                </button>
            </div>
        `).join('');
        document.getElementById('product-list').innerHTML = html || "<p class='text-slate-400'>No food found.</p>";
    } catch (err) { console.error(err); }
}

/**
 * Place Order logic
 */
async function placeOrder(productId) {
    try {
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
    } catch (err) { console.error(err); }
}

/**
 * Fetch Order History
 */
async function loadOrders() {
    try {
        const res = await fetch(`${GATEWAY}/orders`, {
            headers: { 'Authorization': `Bearer ${TOKEN}` }
        });
        if (!res.ok) return;
        const data = await res.json();
        
        document.getElementById('order-count').innerText = data.length;
        document.getElementById('order-history').innerHTML = data.map(o => `
            <div class="bg-slate-800 p-4 rounded-xl border border-slate-700 order-card relative group">
                ${o.status !== 'PENDING' ? `<button onclick="updateStatus('${o.id}', 'ARCHIVED')" class="absolute -top-1 -right-1 bg-red-500 text-white w-5 h-5 rounded-full text-[10px] opacity-0 group-hover:opacity-100 transition">✕</button>` : ''}
                <div class="flex justify-between text-[10px] text-slate-400 mb-1">
                    <span>#${o.id.substring(0,8)}</span>
                    <span class="${o.status === 'COMPLETED' ? 'text-green-400' : 'text-orange-400'} font-bold">${o.status}</span>
                </div>
                <div class="flex justify-between items-center">
                    <span class="font-bold text-lg text-white">$${o.totalPrice}</span>
                    ${o.status === 'PENDING' ? `<button onclick="updateStatus('${o.id}', 'COMPLETED')" class="text-[10px] bg-white text-slate-900 px-2 py-1 rounded-full font-bold hover:bg-orange-100">Mark Delivered</button>` : ''}
                </div>
            </div>
        `).join('');
    } catch (err) { console.error(err); }
}

/**
 * Update Order Status
 */
async function updateStatus(orderId, newStatus) {
    try {
        await fetch(`${GATEWAY}/orders/${orderId}/status`, {
            method: 'PUT',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${TOKEN}`
            },
            body: JSON.stringify({ status: newStatus })
        });
        loadOrders();
    } catch (err) { console.error(err); }
}