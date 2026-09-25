export function setTokens(accessToken, refreshToken) {
    localStorage.setItem('access_token', accessToken);
    localStorage.setItem('refresh_token', refreshToken);
}

export function getAccessToken() {
    return localStorage.getItem('access_token');
}

export function getRefreshToken() {
    return localStorage.getItem('refresh_token');
}

export function clearTokens() {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
}

// Wrapper alrededor de fetch() para manejar autenticación y refresco automático
export async function authFetch(url, options = {}) {
    options.headers = options.headers || {};
    
    let accessToken = getAccessToken();
    if (accessToken) {
        options.headers['Authorization'] = `Bearer ${accessToken}`;
    }

    let response = await fetch(url, options);

    // Si devuelve 401 (token expirado), intentamos hacer refresh silencioso
    if (response.status === 401) {
        const refreshToken = getRefreshToken();
        if (!refreshToken) {
            clearTokens();
            window.location.href = '/login.html';
            return response;
        }

        // Solicitar renovar tokens
        const refreshResponse = await fetch('/refresh', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refresh_token: refreshToken })
        });

        if (refreshResponse.ok) {
            const refreshData = await refreshResponse.json();
            // Guardar los nuevos tokens
            setTokens(refreshData.data.access_token, refreshData.data.refresh_token);

            // Reintentar la petición original con el nuevo token
            options.headers['Authorization'] = `Bearer ${refreshData.data.access_token}`;
            response = await fetch(url, options);
        } else {
            // Si el refresh token también expiró o falló, redirigir al login
            clearTokens();
            window.location.href = '/login.html';
        }
    }

    return response;
}