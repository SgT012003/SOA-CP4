// app.js

const API_BASE = '/api';

function getToken() {
    return localStorage.getItem('token');
}

function authHeaders() {
    return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getToken()}`
    };
}

function logout() {
    localStorage.removeItem('token');
    window.location.href = '/login';
}

// Intercepta e checa autenticacao basica para qualquer pagina interna
if (!getToken() && window.location.pathname !== '/login') {
    window.location.href = '/login';
}

// Se estiver logado e na home, redirecionar para salas
if (getToken() && (window.location.pathname === '/' || window.location.pathname === '/login')) {
    window.location.href = '/salas';
}

async function handleLogin(e) {
    e.preventDefault();
    const email = document.getElementById('login-email').value;
    const senha = document.getElementById('login-senha').value;
    const alertBox = document.getElementById('auth-alert');

    try {
        const res = await fetch(`${API_BASE}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, senha })
        });
        
        const data = await res.json();
        
        if (!res.ok) throw new Error(data.error || 'Erro no login');

        localStorage.setItem('token', data.token);
        window.location.href = '/salas';
    } catch (err) {
        showError(alertBox, err.message);
    }
}

async function handleRegister(e) {
    e.preventDefault();
    const nome = document.getElementById('reg-nome').value;
    const email = document.getElementById('reg-email').value;
    const senha = document.getElementById('reg-senha').value;
    const alertBox = document.getElementById('auth-alert');

    try {
        const res = await fetch(`${API_BASE}/auth/registrar`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ nome, email, senha })
        });
        
        const data = await res.json();
        if (!res.ok) throw new Error(data.error || 'Erro no registro');

        showSuccess(alertBox, 'Conta criada! Faça login.');
        document.getElementById('register-form').reset();
        setTimeout(toggleAuthMode, 1500);
    } catch (err) {
        showError(alertBox, err.message);
    }
}

function showError(el, msg) {
    el.textContent = msg;
    el.className = 'mt-4 p-3 rounded-lg text-sm font-medium text-center bg-red-50 text-red-600 border border-red-200 block';
}

function showSuccess(el, msg) {
    el.textContent = msg;
    el.className = 'mt-4 p-3 rounded-lg text-sm font-medium text-center bg-green-50 text-green-600 border border-green-200 block';
}

// ------- SALAS -------

async function carregarSalas() {
    try {
        const res = await fetch(`${API_BASE}/salas`, { headers: authHeaders() });
        if (res.status === 401) return logout();
        const salas = await res.json();
        
        document.getElementById('loading').classList.add('hidden');
        const grid = document.getElementById('salas-grid');
        grid.innerHTML = '';
        grid.classList.remove('hidden');

        if (salas.length === 0) {
            grid.innerHTML = '<div class="col-span-full text-center text-slate-500 py-8">Nenhuma sala cadastrada.</div>';
            return;
        }

        salas.forEach(s => {
            const statusColor = s.status === 'ATIVA' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700';
            const card = `
                <div class="bg-white rounded-2xl p-6 shadow-sm border border-slate-100 hover:shadow-lg transition-all hover:scale-[1.02]">
                    <div class="flex justify-between items-start mb-4">
                        <h3 class="font-bold text-lg text-slate-800">${s.nome}</h3>
                        <span class="px-2.5 py-1 rounded-full text-xs font-semibold ${statusColor}">${s.status}</span>
                    </div>
                    <div class="text-sm text-slate-600 space-y-2">
                        <p class="flex items-center"><svg class="w-4 h-4 mr-2 text-pink-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path></svg> Capacidade: ${s.capacidade} pessoas</p>
                        <p class="flex items-center"><svg class="w-4 h-4 mr-2 text-cyan-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"></path></svg> ${s.localizacao}</p>
                    </div>
                </div>
            `;
            grid.insertAdjacentHTML('beforeend', card);
        });
    } catch (err) {
        console.error(err);
        alert('Erro ao carregar salas');
    }
}

async function handleNovaSala(e) {
    e.preventDefault();
    const nome = document.getElementById('sala-nome').value;
    const capacidade = parseInt(document.getElementById('sala-capac').value, 10);
    const localizacao = document.getElementById('sala-local').value;
    const status = document.getElementById('sala-status').value;

    try {
        const res = await fetch(`${API_BASE}/salas`, {
            method: 'POST',
            headers: authHeaders(),
            body: JSON.stringify({ nome, capacidade, localizacao, status })
        });
        
        if (!res.ok) {
            const data = await res.json();
            throw new Error(data.error);
        }
        
        document.getElementById('modal-sala').classList.add('hidden');
        document.getElementById('form-sala').reset();
        carregarSalas();
    } catch (err) {
        alert(err.message);
    }
}

// ------- RESERVAS -------

async function carregarReservas() {
    try {
        const res = await fetch(`${API_BASE}/reservas`, { headers: authHeaders() });
        if (res.status === 401) return logout();
        const reservas = await res.json();
        
        const tbody = document.getElementById('tbl-reservas');
        tbody.innerHTML = '';

        if (reservas.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" class="p-8 text-center text-slate-500">Nenhuma reserva encontrada.</td></tr>';
            return;
        }

        reservas.forEach(r => {
            // "2024-05-15T00:00:00Z" -> "15/05/2024"
            const dataStr = new Date(r.data).toLocaleDateString('pt-BR', { timeZone: 'UTC' });
            // "09:00:00" -> "09:00"
            const hI = r.horario_inicio.substring(0, 5);
            const hF = r.horario_fim.substring(0, 5);

            const statusUI = r.status === 'CONFIRMADA' 
                ? '<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">Confirmada</span>'
                : '<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">Cancelada</span>';

            const btnCancel = r.status === 'CONFIRMADA' ? `<button onclick="cancelarReserva('${r.id}')" class="text-red-500 hover:text-red-700 font-medium text-sm transition-colors">Cancelar</button>` : `<span class="text-slate-300">-</span>`;

            const tr = `
                <tr class="hover:bg-slate-50/50 transition-colors">
                    <td class="p-4"><div class="font-medium text-slate-800">${r.nome_sala}</div></td>
                    <td class="p-4 hidden sm:table-cell text-slate-600">${r.finalidade || '-'}</td>
                    <td class="p-4 text-slate-600">
                        <div class="font-medium">${dataStr}</div>
                        <div class="text-xs text-slate-400">${hI} - ${hF}</div>
                    </td>
                    <td class="p-4 hidden md:table-cell text-slate-600">
                        <div class="flex items-center">
                            <div class="w-6 h-6 rounded-full bg-cyan-100 text-cyan-600 flex items-center justify-center text-xs font-bold mr-2">
                                ${r.nome_usuario.charAt(0).toUpperCase()}
                            </div>
                            ${r.nome_usuario}
                        </div>
                    </td>
                    <td class="p-4 text-center">${statusUI}</td>
                    <td class="p-4 text-right">${btnCancel}</td>
                </tr>
            `;
            tbody.insertAdjacentHTML('beforeend', tr);
        });
    } catch (err) {
        console.error(err);
    }
}

async function abrirModalReserva() {
    document.getElementById('modal-reserva').classList.remove('hidden');
    const select = document.getElementById('reserva-sala');
    select.innerHTML = '<option value="">Carregando salas...</option>';
    
    try {
        const res = await fetch(`${API_BASE}/salas`, { headers: authHeaders() });
        const salas = await res.json();
        const ativas = salas.filter(s => s.status === 'ATIVA');
        if (ativas.length === 0) {
            select.innerHTML = '<option value="">Nenhuma sala ATIVA disponível</option>';
            return;
        }
        select.innerHTML = '<option value="" disabled selected>Selecione uma sala</option>';
        ativas.forEach(s => {
            select.innerHTML += `<option value="${s.id}">${s.nome} (${s.capacidade} p.)</option>`;
        });
    } catch (err) {
        select.innerHTML = '<option value="">Erro ao carregar</option>';
    }
}

async function handleNovaReserva(e) {
    e.preventDefault();
    const id_sala = document.getElementById('reserva-sala').value;
    const data = document.getElementById('reserva-data').value;
    // Precisamos enviar HH:MM:00 ou o banco aceita HH:MM. Postgres Time aceita HH:MM.
    const horario_inicio = document.getElementById('reserva-inicio').value;
    const horario_fim = document.getElementById('reserva-fim').value;
    const finalidade = document.getElementById('reserva-motivo').value;

    if (!id_sala) return alert('Selecione uma sala válida');

    try {
        const res = await fetch(`${API_BASE}/reservas`, {
            method: 'POST',
            headers: authHeaders(),
            body: JSON.stringify({ id_sala, data, horario_inicio, horario_fim, finalidade })
        });
        
        if (!res.ok) {
            const errData = await res.json();
            throw new Error(errData.error || 'Erro ao criar reserva');
        }
        
        document.getElementById('modal-reserva').classList.add('hidden');
        document.getElementById('form-reserva').reset();
        carregarReservas();
    } catch (err) {
        alert(err.message);
    }
}

async function cancelarReserva(id) {
    if (!confirm('Tem certeza que deseja cancelar esta reserva?')) return;

    try {
        const res = await fetch(`${API_BASE}/reservas/${id}/cancelar`, {
            method: 'PATCH',
            headers: authHeaders()
        });
        
        if (!res.ok) {
            const errData = await res.json();
            throw new Error(errData.error || 'Erro ao cancelar');
        }
        
        carregarReservas();
    } catch (err) {
        alert(err.message);
    }
}
