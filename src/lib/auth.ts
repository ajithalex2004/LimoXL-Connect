import { create } from 'zustand';
import axios from 'axios';

const isProd = import.meta.env.PROD;
const apiPrefix = import.meta.env.VITE_API_URL || (isProd ? '/api' : 'http://localhost:8080/api');

// Axios Instance
export const api = axios.create({
    baseURL: apiPrefix,
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

// Helper to safely parse JSON from localStorage
const safeParse = (key: string, fallback: any) => {
    try {
        const item = localStorage.getItem(key);
        if (item === null || item === 'undefined') return fallback;
        return JSON.parse(item);
    } catch (e) {
        console.error(`Error parsing localStorage key "${key}":`, e);
        return fallback;
    }
};

// Auth Store
export const useAuthStore = create<AuthState>((set) => {
    const token = localStorage.getItem('token');
    const user = safeParse('user', null);

    return {
        token,
        user,
        isAuthenticated: !!token && !!user,

        login: (token, user) => {
            console.log('DEBUG: Store login called with:', { token: !!token, user });
            if (!token || !user) {
                console.error('DEBUG: login called with missing data, aborting save');
                return;
            }
            localStorage.setItem('token', token);
            localStorage.setItem('user', JSON.stringify(user));
            api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
            set({ token, user, isAuthenticated: true });
        },

    },

        logout: () => {
            console.log('DEBUG: Store logout called');
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            delete api.defaults.headers.common['Authorization'];
            set({ token: null, user: null, isAuthenticated: false });
        },
    };
});

// Add interceptor to attach token
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
            console.log('DEBUG: Storage event triggered for key:', e.key);
            // Another tab changed the auth state
            const newToken = localStorage.getItem('token');
            const newUser = safeParse('user', null);

            // Update the store to reflect the change
            useAuthStore.setState({
                token: newToken,
                user: newUser,
                isAuthenticated: !!newToken
            });

            // Optionally reload the page to prevent confusion
            if (e.oldValue !== e.newValue) {
                console.warn('Auth state changed in another tab. Reloading...');
                // window.location.reload(); // Temporarily disabled to avoid lost logs during debug
            }
        }
    });
}
