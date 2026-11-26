import { setTokens, getRefreshToken, clearTokens, authFetch } from './api.js';
// dashboar INICIO

// --- VALIDACIÓN DE ACCESO AL DASHBOARD ---
// Verificamos si estamos en la página del dashboard
if (window.location.pathname === '/dashboard') {
    const accessToken = localStorage.getItem('access_token'); // o sessionStorage según donde los guardes

    // Si no hay token, fuera del dashboard inmediatamente
    if (!accessToken) {
        window.location.href = '/login';
    } else {
        // Si hay token, cargamos los datos del usuario para el topbar
        cargarDatosPerfil(accessToken);
    }
}

// Función para poblar el nombre y correo en la esquina superior izquierda/derecha
// --- 1. FUNCIÓN PARA CARGAR EL PERFIL ---
async function cargarDatosPerfil(token) {
    try {
        const res = await fetch('/api/profile', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}` // El token viaja al endpoint protegido de Go
            }
        });

        const data = await res.json();

        if (res.ok) {
            // Rellenamos los elementos de la topbar (asegúrate de que existan en el HTML)
            const userNameEl = document.getElementById('userName');
            const userEmailEl = document.getElementById('userEmail');

            if (userNameEl) userNameEl.textContent = data.data.nombre || 'Usuario';
            if (userEmailEl) userEmailEl.textContent = data.data.correo || '';
        } else {
            // Si el token expiró o es inválido, limpiamos y mandamos al login
            localStorage.clear();
            window.location.href = '/login';
        }
    } catch (err) {
        console.error('Error al obtener perfil', err);
    }
}

// --- 2. DISPARADOR AUTOMÁTICO AL CARGAR CUALQUIER PÁGINA ---
document.addEventListener('DOMContentLoaded', () => {
    const token = localStorage.getItem('access_token');
    
    // Si estamos en una página protegida y no hay token, fuera al login
    if (!token && window.location.pathname !== '/login' && window.location.pathname !== '/register') {
        window.location.href = '/login';
        return;
    }

    // Si hay token, cargamos los datos del perfil automáticamente en la topbar
    if (token) {
        cargarDatosPerfil(token);
    }
});



// dashboard FIN

// siderbar

// --- CONTROL DEL SIDEBAR COLAPSIBLE ---
const sidebarToggle = document.getElementById('sidebarToggle');
const sidebar = document.querySelector('.sidebar');

if (sidebarToggle && sidebar) {
    sidebarToggle.addEventListener('click', () => {
        sidebar.classList.toggle('collapsed');
    });
}

// --- CONTROL DE MENÚS DESPLEGABLES (GRUPOS) ---
const dropdownToggles = document.querySelectorAll('.dropdown-toggle');

dropdownToggles.forEach(toggle => {
    toggle.addEventListener('click', (e) => {
        e.preventDefault();
        const parentLi = toggle.parentElement;
        
        // Alterna la clase 'open' para expandir/contraer el submenú
        parentLi.classList.toggle('open');
    });
});

// siderbar fin

//luego de crear usuario, en actualizar password, el usuario no volverá al colocar su correo de forma manual
const firstEmailInput = document.getElementById('firstEmail');
const savedEmail = sessionStorage.getItem('first_login_email');

if (firstEmailInput && savedEmail) {
    firstEmailInput.value = savedEmail;
    firstEmailInput.readOnly = true; 
}

//

document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('loginForm');
    const btnLogout = document.getElementById('btnLogout');
    const profileData = document.getElementById('profileData');

    // lógica para el login
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const correo = document.getElementById('correo').value;
            const clave = document.getElementById('clave').value;
            const errorMsg = document.getElementById('errorMsg');

            errorMsg.textContent = '';

            try {
                const res = await fetch('/api/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ correo, clave })
                });

                const data = await res.json();

                if (res.ok) {
                    // pedir configurar password cuando es la primera vez que ingresa
                    if (data.data && data.data.requires_password_setup) {
                        sessionStorage.setItem('first_login_email', data.data.correo);
                        window.location.href = '/update-password';
                        return; // Detiene la ejecución para no guardar tokens ni ir a /profile
                    }

                    // flujo normal cuando ya tiene contraseña configurada
                    setTokens(data.data.access_token, data.data.refresh_token);
                    window.location.href = '/dashboard';
                } else {
                    errorMsg.textContent = data.error || 'Credenciales inválidas';
                }
            } catch (err) {
                errorMsg.textContent = 'Error de conexión con el servidor';
            }
        });
    }

    // vista del profile
    if (profileData) {
        loadProfile();
    }

    // logout
    if (btnLogout) {
        btnLogout.addEventListener('click', async () => {
            const refreshToken = getRefreshToken();
            if (refreshToken) {
                await authFetch('/logout', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ refresh_token: refreshToken })
                });
            }
            clearTokens();
            window.location.href = '/login';
        });
    }
});

async function loadProfile() {
    const profileData = document.getElementById('profileData');
    try {
        const res = await authFetch('/api/profile');
        if (res.ok) {
            const data = await res.json();
            const p = data.data;
            profileData.innerHTML = `
                <p><strong>Nombre:</strong> ${p.apellido_paterno_registrador} ${p.nombres_registrador}</p>
                <p><strong>Correo:</strong> ${p.correo}</p>
            `;
        }
    } catch (err) {
        profileData.innerHTML = `<p class="error">Error al cargar datos del perfil</p>`;
    }
}

// Para el registro de usuarios
const registerForm = document.getElementById('registerForm');

if (registerForm) {
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const documentNumber = document.getElementById('documentNumber').value.trim();
        const email = document.getElementById('regCorreo').value.trim();
        const regMsg = document.getElementById('regMsg');

        regMsg.textContent = '';
        regMsg.style.color = '#d9534f'; // Color rojo predeterminado para errores

        try {
            const res = await fetch('/api/CreateUsers', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    numero_documento: documentNumber,
                    correo: email
                })
            });

            const data = await res.json();

            if (res.ok) {
                // Mensaje de éxito
                regMsg.style.color = '#28a745'; // Verde
                regMsg.textContent = 'Usuario registrado exitosamente. Redirigiendo al login...';

                // Redirigir al inicio de sesión tras 2 segundos
                setTimeout(() => {
                    window.location.href = '/login';
                }, 2000);
            } else {
                regMsg.textContent = data.error || data.err || 'Error al registrar usuario';
            }
        } catch (err) {
            regMsg.textContent = 'Error de conexión con el servidor';
        }
    });
}

// primer inicio de sesión
const firstLoginForm = document.getElementById('firstLoginForm');

if (firstLoginForm) {
    firstLoginForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const email = document.getElementById('firstEmail').value.trim();
        const password = document.getElementById('firstPassword').value;
        const confirmPassword = document.getElementById('firstConfirmPassword').value;
        const msgElement = document.getElementById('firstLoginMsg');

        msgElement.textContent = '';
        msgElement.style.color = '#d9534f'; // Rojo para errores

        // Validación preventiva en cliente
        if (password !== confirmPassword) {
            msgElement.textContent = 'Las contraseñas no coinciden';
            return;
        }

        try {
            const res = await fetch('/api/UpdatePassword', {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    correo: email,
                    clave: password,
                    confirmar_clave: confirmPassword
                })
            });

            const data = await res.json();

            if (res.ok) {
                msgElement.style.color = '#28a745'; // Verde
                msgElement.textContent = 'Contraseña establecida exitosamente. Redirigiendo al login...';

                setTimeout(() => {
                    window.location.href = '/login';
                }, 2000);
            } else {
                msgElement.textContent = data.error || data.err || 'Error al actualizar contraseña';
            }
        } catch (err) {
            msgElement.textContent = 'Error de conexión con el servidor';
        }
    });
}

// para el reseteo de password
const resetPasswordForm = document.getElementById('resetPasswordForm');
if(resetPasswordForm) {
    resetPasswordForm.addEventListener('submit', async (e) =>{
        e.preventDefault();

        const documentNumber = document.getElementById('documentNumber').value.trim();
        const email = document.getElementById('resetCorreo').value;
        const newPassword = document.getElementById('newPassword').value;
        const newConfirmPassword = document.getElementById('newConfirmPassword').value;

        const msgElement = document.getElementById('resetMsg')
        msgElement.textContent = '';
        msgElement.style.color = '#d9534f'; // rojo para errores

        // validacion previa en cliente
        if (newPassword != newConfirmPassword) {
            msgElement.textContent = 'Las contraseñas no coinciden';
            return;
        }

        try {
            const res = await fetch('/api/ResetPassword', {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    numero_documento: documentNumber,
                    correo: email,
                    new_password: newPassword,
                    confirm_password: newConfirmPassword
                })
            });

            const data = await res.json();

            if (res.ok) {
                msgElement.style.color = '#28a745'
                msgElement.textContent = 'Constraseña reestablecida con éxito. Redirigiendo al login...';

                setTimeout(() => {
                    window.location.href = '/login';
                }, 2000);
            } else {
                msgElement.textContent = data.error || data.err || 'Error al actualizar contraseña';
            }
        } catch (err) {
            msgElement.textContent = 'Error de conexión con el servidor'
        }
    });
}

// --- LÓGICA DE SELECTS EN CASCADA (REPORTE VERTICAL Y BÚSQUEDA ESPECÍFICA) ---
const esReporteVertical = window.location.pathname === '/vertical-report';
const esBusquedaEspecifica = window.location.pathname === '/filter-report';

if (esReporteVertical || esBusquedaEspecifica) {
    const selectRed = document.getElementById('selectRed');
    const selectMicrored = document.getElementById('selectMicrored');
    const selectPd = document.getElementById('selectPd');
    const selectEstablecimiento = document.getElementById('selectEstablecimiento');
    const selectCursoVida = document.getElementById('selectCursoVida');
    const token = localStorage.getItem('access_token');

    // 1. Cargar Redes al iniciar la página
    async function cargarRedes() {
        if (!selectRed) return;
        try {
            const res = await fetch('/api/reportes/redes', {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const json = await res.json();
            if (res.ok && Array.isArray(json.data)) {
                json.data.forEach(red => {
                    const opt = document.createElement('option');
                    opt.value = red;
                    opt.textContent = red;
                    selectRed.appendChild(opt);
                });
            }
        } catch (err) {
            console.error('Error cargando redes:', err);
        }
    }

    // Evento Red -> Carga Microredes
    if (selectRed) {
        selectRed.addEventListener('change', async () => {
            const redVal = selectRed.value;

            selectMicrored.innerHTML = '<option value="">-- Seleccione Microred --</option>';
            selectMicrored.disabled = !redVal;
            selectPd.innerHTML = '<option value="">-- Seleccione P. Digitacion --</option>';
            selectPd.disabled = true;
            selectEstablecimiento.innerHTML = '<option value="">-- Seleccione Establecimiento --</option>';
            selectEstablecimiento.disabled = true;

            if (!redVal) return;

            try {
                const res = await fetch(`/api/reportes/microredes?red=${encodeURIComponent(redVal)}`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                const json = await res.json();
                if (res.ok && Array.isArray(json.data)) {
                    json.data.forEach(mr => {
                        const opt = document.createElement('option');
                        opt.value = mr;
                        opt.textContent = mr;
                        selectMicrored.appendChild(opt);
                    });
                }
            } catch (err) {
                console.error('Error cargando microredes:', err);
            }
        });
    }

    // Evento Microred -> Carga Puntos de Digitación
    if (selectMicrored) {
        selectMicrored.addEventListener('change', async () => {
            const redVal = selectRed.value;
            const mrVal = selectMicrored.value;

            selectPd.innerHTML = '<option value="">-- Seleccione P. Digitacion --</option>';
            selectPd.disabled = !mrVal;
            selectEstablecimiento.innerHTML = '<option value="">-- Seleccione Establecimiento --</option>';
            selectEstablecimiento.disabled = true;

            if (!mrVal) return;

            try {
                const res = await fetch(`/api/reportes/pd-digitacion?red=${encodeURIComponent(redVal)}&microred=${encodeURIComponent(mrVal)}`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                const json = await res.json();
                if (res.ok && Array.isArray(json.data)) {
                    json.data.forEach(pd => {
                        const opt = document.createElement('option');
                        opt.value = pd;
                        opt.textContent = pd;
                        selectPd.appendChild(opt);
                    });
                }
            } catch (err) {
                console.error('Error cargando puntos de digitación:', err);
            }
        });
    }

    // Evento P. Digitacion -> Carga Establecimientos
    if (selectPd) {
        selectPd.addEventListener('change', async () => {
            const redVal = selectRed.value;
            const mrVal = selectMicrored.value;
            const pdVal = selectPd.value;

            selectEstablecimiento.innerHTML = '<option value="">-- Seleccione Establecimiento --</option>';
            selectEstablecimiento.disabled = !pdVal;

            if (!pdVal) return;

            try {
                const res = await fetch(`/api/reportes/establecimientos?red=${encodeURIComponent(redVal)}&microred=${encodeURIComponent(mrVal)}&pd_digitacion=${encodeURIComponent(pdVal)}`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                const json = await res.json();
                if (res.ok && Array.isArray(json.data)) {
                    json.data.forEach(est => {
                        const opt = document.createElement('option');
                        opt.value = est;
                        opt.textContent = est;
                        selectEstablecimiento.appendChild(opt);
                    });
                }
            } catch (err) {
                console.error('Error cargando establecimientos:', err);
            }
        });
    }

    // Cargar Cursos de Vida solo si existe el select (en Reporte Vertical)
    async function cargarCursosDeVida() {
        if (!selectCursoVida) return;

        try {
            const res = await fetch('/api/reportes/curso-vida', {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const json = await res.json();

            if (res.ok) {
                const lista = json.data || json;
                if (Array.isArray(lista)) {
                    lista.forEach(item => {
                        const opt = document.createElement('option');
                        opt.value = item;
                        opt.textContent = item;
                        selectCursoVida.appendChild(opt);
                    });
                }
            }
        } catch (err) {
            console.error('Error cargando cursos de vida:', err);
        }
    }

    cargarRedes();
    if (esReporteVertical) {
        cargarCursosDeVida();
    }
}

// --- LÓGICA DE BÚSQUEDA, PAGINACIÓN Y EXPORTACIÓN ---
const filterForm = document.getElementById('filterForm');
const resultsBody = document.getElementById('resultsBody');
const btnPrev = document.getElementById('btnPrev');
const btnNext = document.getElementById('btnNext');
const pageInfo = document.getElementById('pageInfo');
const btnReset = document.getElementById('btnReset');
const btnExportar = document.getElementById('btnExportar');

let currentPage = 1;
const limit = esBusquedaEspecifica ? 22 : 10;

function construirPayload() {
    const basePayload = {
        anio_p: parseInt(document.getElementById('inputAnio')?.value, 10) || 0,
        mes_inicio: parseInt(document.getElementById('inputMesInicio')?.value, 10) || 0,
        mes_final: parseInt(document.getElementById('inputMesFinal')?.value, 10) || 0,
        red_p: document.getElementById('selectRed')?.value || "",
        microred_p: document.getElementById('selectMicrored')?.value || "",
        pd_digitacion_p: document.getElementById('selectPd')?.value || "",
        establecimiento_p: document.getElementById('selectEstablecimiento')?.value || "",
        documento_p: document.getElementById('inputDocumento')?.value.trim() || "",
        cie_cmp_p: document.getElementById('inputCieCpm')?.value.trim() || ""
    };

    if (esBusquedaEspecifica) {
        return {
            ...basePayload,
            id_cita_p: parseInt(document.getElementById('inputIdCita')?.value, 10) || 0,
            app_p: document.getElementById('inputApp')?.value.trim() || "",
            lab_p: document.getElementById('inputLab')?.value.trim() || ""
        };
    }

    return {
        ...basePayload,
        curso_vida_p: document.getElementById('selectCursoVida')?.value || ""
    };
}

function renderizarFilas(registros) {
    resultsBody.innerHTML = '';

    registros.forEach(item => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>${item.id_cita || ''}</td>
            <td>${item.anio || ''}</td>
            <td>${item.mes || ''}</td>
            <td>${item.dia || ''}</td>
            <td>${item.lote || ''}</td>
            <td>${item.num_pag || ''}</td>
            <td>${item.num_reg || ''}</td>
            <td>${item.codigo_unico || ''}</td>
            <td>${item.nombre_establecimiento || ''}</td>
            <td>${item.id_paciente || ''}</td>
            <td>${item.documento || ''}</td>
            <td>${item.nombre_paciente || ''}</td>
            <td>${item.genero || ''}</td>
            <td>${item.id_condicion_establecimiento || ''}</td>
            <td>${item.id_condicion_servicio || ''}</td>
            <td>${item.edad_reg || ''}</td>
            <td>${item.tipo_edad || ''}</td>
            <td>${item.codigo_item || ''}</td>
            <td>${item.tipo_diagnostico || ''}</td>
            <td>${item.valor_lab || ''}</td>
            <td>${item.id_correlativo || ''}</td>
            <td>${item.id_correlativo_lab || ''}</td>
            <td>${item.peso || ''}</td>
            <td>${item.talla || ''}</td>
            <td>${item.hemoglobina || ''}</td>
            <td>${item.perimetro_abdominal || ''}</td>
            <td>${item.perimetro_cefalico || ''}</td>
            <td>${item.descripcion_otra_condicion || ''}</td>
            <td>${item.fecha_ultima_regla || ''}</td>
            <td>${item.nombre_personal_salud || ''}</td>
            <td>${item.nombre_registrador || ''}</td>
            <td>${item.fecha_registro || ''}</td>
            <td>${item.fecha_modificacion || ''}</td>
        `;
        resultsBody.appendChild(tr);
    });
}

async function realizarBusqueda(page = 1) {
    currentPage = page;
    const payload = construirPayload();

    console.log(">>> PAYLOAD ENVIADO A GO (Pág " + currentPage + "):", JSON.stringify(payload, null, 2));

    resultsBody.innerHTML = `<tr><td colspan="33" class="text-center">Buscando registros...</td></tr>`;

    try {
        const token = localStorage.getItem('access_token');
        if (!token) {
            alert("Tu sesión ha expirado. Por favor, inicia sesión nuevamente.");
            window.location.href = '/login';
            return;
        }

        const endpointBuscar = esBusquedaEspecifica
            ? `/api/buscaespecifica?page=${currentPage}&limit=${limit}`
            : `/api/buscar?page=${currentPage}&limit=${limit}`;

        const res = await fetch(endpointBuscar, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(payload)
        });

        const json = await res.json();
        console.log(">>> RESPUESTA DEL SERVIDOR:", json);

        if (res.ok) {
            const registros = json.data || [];

            if (!Array.isArray(registros) || registros.length === 0) {
                resultsBody.innerHTML = `<tr><td colspan="33" class="text-center">No se encontraron registros con los filtros seleccionados.</td></tr>`;
                if (pageInfo) pageInfo.textContent = `Página 0 de 0`;
                if (btnPrev) btnPrev.disabled = true;
                if (btnNext) btnNext.disabled = true;
                return;
            }

            renderizarFilas(registros);

            const meta = json.meta || { page: 1, page_count: 1 };
            if (pageInfo) pageInfo.textContent = `Página ${meta.page} de ${meta.page_count}`;
            if (btnPrev) btnPrev.disabled = meta.page <= 1;
            if (btnNext) btnNext.disabled = meta.page >= meta.page_count;
        } else {
            console.error('Error del servidor:', json);
            resultsBody.innerHTML = `<tr><td colspan="33" class="text-center" style="color: red;">Error en la respuesta del servidor.</td></tr>`;
        }
    } catch (err) {
        console.error('Error de red en la búsqueda:', err);
        resultsBody.innerHTML = `<tr><td colspan="33" class="text-center" style="color: red;">Error de conexión con el servidor.</td></tr>`;
    }
}

if (filterForm) {
    filterForm.addEventListener('submit', (e) => {
        e.preventDefault();
        realizarBusqueda(1);
    });
}

if (btnPrev) {
    btnPrev.addEventListener('click', () => {
        if (currentPage > 1) {
            realizarBusqueda(currentPage - 1);
        }
    });
}

if (btnNext) {
    btnNext.addEventListener('click', () => {
        realizarBusqueda(currentPage + 1);
    });
}

if (btnReset && filterForm) {
    btnReset.addEventListener('click', () => {
        filterForm.reset();
        const selectMicrored = document.getElementById('selectMicrored');
        const selectPd = document.getElementById('selectPd');
        const selectEstablecimiento = document.getElementById('selectEstablecimiento');

        if (selectMicrored) {
            selectMicrored.innerHTML = '<option value="">-- Seleccione Microred --</option>';
            selectMicrored.disabled = true;
        }
        if (selectPd) {
            selectPd.innerHTML = '<option value="">-- Seleccione P. Digitacion --</option>';
            selectPd.disabled = true;
        }
        if (selectEstablecimiento) {
            selectEstablecimiento.innerHTML = '<option value="">-- Seleccione Establecimiento --</option>';
            selectEstablecimiento.disabled = true;
        }
    });
}

if (btnExportar) {
    btnExportar.addEventListener('click', async () => {
        const payload = construirPayload();
        const token = localStorage.getItem('access_token');

        if (!token) {
            alert("Tu sesión ha expirado.");
            window.location.href = '/login';
            return;
        }

        btnExportar.textContent = "Generando Excel...";
        btnExportar.disabled = true;

        try {
            const endpointExportar = esBusquedaEspecifica
                ? '/api/reportes/exportarbusca'
                : '/api/reportes/exportar';

            const res = await fetch(endpointExportar, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(payload)
            });

            if (!res.ok) {
                throw new Error("Error al generar el archivo Excel");
            }

            const blob = await res.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            const prefijo = esBusquedaEspecifica ? 'Busqueda_Especifica' : 'Reporte_Vertical';
            a.href = url;
            a.download = `${prefijo}_${new Date().toISOString().slice(0, 10)}.xlsx`;
            document.body.appendChild(a);
            a.click();
            a.remove();
            window.URL.revokeObjectURL(url);
        } catch (err) {
            console.error("Error en la exportación:", err);
            alert("Ocurrió un error al exportar los datos.");
        } finally {
            btnExportar.textContent = "Exportar a Excel";
            btnExportar.disabled = false;
        }
    });
}