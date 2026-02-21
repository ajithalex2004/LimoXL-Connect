import { create } from 'zustand';
import axios from 'axios';

// Axios Instance
export const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',
    headers: {
        'Content-Type': 'application/json',
    },
});

interface User {
    id: string; // Needed for update
    email: string;
    name: string;
    role: string;
    password_change_required?: boolean;
}

interface AuthState {
    token: string | null;
    user: User | null;
    isAuthenticated: boolean;
    login: (token: string, user: User) => void;
    logout: () => void;
}

// Auth Store
export const useAuthStore = create<AuthState>((set) => ({
    token: localStorage.getItem('token'),
    user: JSON.parse(localStorage.getItem('user') || 'null'),
    isAuthenticated: !!localStorage.getItem('token'),

    login: (token, user) => {
        localStorage.setItem('token', token);
        localStorage.setItem('user', JSON.stringify(user));
        api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
        set({ token, user, isAuthenticated: true });
    },

    logout: () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        delete api.defaults.headers.common['Authorization'];
        set({ token: null, user: null, isAuthenticated: false });
    },
}));

// Add interceptor to attach token to requests
api.interceptors.request.use((config) => {
    const token = useAuthStore.getState().token;
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// Listen for storage changes from other tabs
if (typeof window !== 'undefined') {
    window.addEventListener('storage', (e) => {
        if (e.key === 'user' || e.key === 'token') {
            // Another tab changed the auth state
            const newToken = localStorage.getItem('token');
            const newUser = JSON.parse(localStorage.getItem('user') || 'null');

            // Update the store to reflect the change
            useAuthStore.setState({
                token: newToken,
                user: newUser,
                isAuthenticated: !!newToken
            });

            // Optionally reload the page to prevent confusion
            if (e.oldValue !== e.newValue) {
                console.warn('Auth state changed in another tab. Reloading...');
                window.location.reload();
            }
        }
    });
}
