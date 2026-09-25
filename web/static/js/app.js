import { setTokens, getRefreshToken, clearTokens, authFetch } from './api.js';

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
                    window.location.href = '/profile';
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