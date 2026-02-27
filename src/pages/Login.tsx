import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore, api } from '../lib/auth';
import { Info, Loader2 } from 'lucide-react';

const Login = () => {
    const [email, setEmail] = useState('admin@limoxlink.com');
    const [password, setPassword] = useState('admin123');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const navigate = useNavigate();
    const { login } = useAuthStore();

    React.useEffect(() => {
        console.log('DEBUG: Axios baseURL check:', api.defaults.baseURL);
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        try {
            console.log('DEBUG: Attempting login for:', email);
            const response = await api.post('auth/login', { email, password });
            console.log('DEBUG: Login response data:', response.data);

            const token = response.data.token || response.data.Token;
            const userData = response.data.user || response.data.User;

            if (!userData) {
                console.error('DEBUG: No user data found in response!', response.data);
                throw new Error('User data missing from server response');
            }

            login(token, userData);
            navigate('/');
        } catch (err: any) {
            console.error(err);
            setError(err.response?.data || 'Failed to login. Please check your credentials.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen bg-slate-900 flex items-center justify-center p-4">
            <div className="bg-slate-800 border border-slate-700 rounded-xl shadow-2xl w-full max-w-md p-8">
                <div className="text-center mb-8">
                    <h1 className="text-3xl font-bold bg-gradient-to-r from-blue-400 to-indigo-500 text-transparent bg-clip-text mb-2">
                        Limo XL Connect
                    </h1>
                    <p className="text-slate-400">Sign in to your account</p>
                </div>

                {error && (
                    <div className="bg-red-500/10 border border-red-500/20 text-red-400 p-3 rounded-lg mb-6 text-sm flex items-center gap-2">
                        <Info size={16} />
                        {error}
                    </div>
                )}

                <form onSubmit={handleSubmit} className="space-y-6">
                    <div className="space-y-2">
                        <label className="text-sm font-medium text-slate-300">Email Address</label>
                        <input
                            type="email"
                            required
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-3 text-white focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all placeholder:text-slate-600"
                            placeholder="admin@limoxlink.com"
                        />
                    </div>

                    <div className="space-y-2">
                        <label className="text-sm font-medium text-slate-300">Password</label>
                        <input
                            type="password"
                            required
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-3 text-white focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all placeholder:text-slate-600"
                            placeholder="••••••••"
                        />
                    </div>

                    <button
                        type="submit"
                        disabled={loading}
                        className="w-full bg-blue-600 hover:bg-blue-500 text-white font-medium py-3 rounded-lg transition-colors flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {loading && <Loader2 size={18} className="animate-spin" />}
                        {loading ? 'Signing in...' : 'Sign In (v2.2)'}
                    </button>

                    <div className="text-center text-xs text-slate-500 mt-4">
                        <p>Demo Credentials:</p>
                        <p>admin@limoxlink.com / admin123</p>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default Login;
