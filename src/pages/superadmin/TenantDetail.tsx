import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { superAdminService, TenantWithFeatures } from '../../services/superadmin';
import { ArrowLeft, Building2, Calendar, Globe, Shield, Users as UsersIcon, Zap, Key, Mail, X } from 'lucide-react';
import { format } from 'date-fns';

const TenantDetail = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [tenant, setTenant] = useState<TenantWithFeatures | null>(null);
    const [loading, setLoading] = useState(true);
    const [isProvisionModalOpen, setIsProvisionModalOpen] = useState(false);
    const [adminName, setAdminName] = useState('');
    const [adminEmail, setAdminEmail] = useState('');
    const [adminPassword, setAdminPassword] = useState('');
    const [isSubmitting, setIsSubmitting] = useState(false);

    useEffect(() => {
        if (id) fetchTenant();
    }, [id]);

    const fetchTenant = async () => {
        try {
            setLoading(true);
            const tenants = await superAdminService.listTenants();
            const found = tenants.find(t => t.id === id);
            if (found) setTenant(found);
        } catch (error) {
            console.error('Failed to fetch tenant:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleProvisionAdmin = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!id) return;

        try {
            setIsSubmitting(true);
            await superAdminService.createTenantAdmin(id, {
                name: adminName,
                email: adminEmail,
                password: adminPassword
            });
            alert('Admin provisioned successfully!');
            setIsProvisionModalOpen(false);
            setAdminName('');
            setAdminEmail('');
            setAdminPassword('');
            fetchTenant();
        } catch (error) {
            console.error('Failed to provision admin:', error);
            alert('Failed to provision admin user');
        } finally {
            setIsSubmitting(false);
        }
    };

    if (loading) return <div className="p-8">Loading...</div>;
    if (!tenant) return <div className="p-8 text-red-500">Tenant not found</div>;

    return (
        <div className="p-6 max-w-5xl mx-auto">
            <button 
                onClick={() => navigate('/superadmin/dashboard')}
                className="flex items-center gap-2 text-gray-500 hover:text-indigo-600 mb-6 transition-colors"
            >
                <ArrowLeft size={18} />
                Back to Orchestration
            </button>

            <div className="bg-white rounded-3xl shadow-sm border border-gray-100 p-8 mb-8">
                <div className="flex items-center gap-6 mb-8">
                    <div className="w-20 h-20 rounded-2xl bg-indigo-600 text-white flex items-center justify-center text-3xl font-bold shadow-lg shadow-indigo-200">
                        {tenant.name.charAt(0)}
                    </div>
                    <div>
                        <h1 className="text-3xl font-bold text-gray-900">{tenant.name}</h1>
                        <div className="flex items-center gap-4 mt-2">
                            <span className="flex items-center gap-1.5 text-sm text-gray-400 font-mono tracking-tighter uppercase">
                                <Globe size={14} />
                                {tenant.slug}
                            </span>
                            <span className="px-3 py-1 bg-emerald-50 text-emerald-700 rounded-full text-xs font-bold border border-emerald-100">
                                {tenant.status}
                            </span>
                        </div>
                    </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                    <div className="flex items-center gap-4 p-4 bg-gray-50 rounded-2xl border border-gray-100">
                        <div className="p-3 bg-white rounded-xl text-indigo-600 shadow-sm">
                            <Shield size={20} />
                        </div>
                        <div>
                            <p className="text-xs text-gray-400 font-bold uppercase tracking-wider">Subscription</p>
                            <p className="text-lg font-bold text-gray-900">{tenant.plan}</p>
                        </div>
                    </div>
                    <div className="flex items-center gap-4 p-4 bg-gray-50 rounded-2xl border border-gray-100">
                        <div className="p-3 bg-white rounded-xl text-amber-600 shadow-sm">
                            <UsersIcon size={20} />
                        </div>
                        <div>
                            <p className="text-xs text-gray-400 font-bold uppercase tracking-wider">User Allocation</p>
                            <p className="text-lg font-bold text-gray-900">{tenant.user_count} / {tenant.max_users}</p>
                        </div>
                    </div>
                    <div className="flex items-center gap-4 p-4 bg-gray-50 rounded-2xl border border-gray-100">
                        <div className="p-3 bg-white rounded-xl text-emerald-600 shadow-sm">
                            <Zap size={20} />
                        </div>
                        <div>
                            <p className="text-xs text-gray-400 font-bold uppercase tracking-wider">Active Instances</p>
                            <p className="text-lg font-bold text-gray-900">{tenant.features.filter(f => f.is_enabled).length} Microservices</p>
                        </div>
                    </div>
                </div>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                <div className="bg-white rounded-3xl shadow-sm border border-gray-100 p-8">
                    <h3 className="text-xl font-bold text-gray-900 mb-6 flex items-center gap-2">
                        <Building2 size={20} className="text-indigo-600" />
                        Infrastructure Info
                    </h3>
                    <div className="space-y-4 text-sm">
                        <div className="flex justify-between py-3 border-b border-gray-50">
                            <span className="text-gray-400">Created At</span>
                            <span className="font-medium text-gray-700">{format(new Date(tenant.created_at), 'PPP')}</span>
                        </div>
                        <div className="flex justify-between py-3 border-b border-gray-50">
                            <span className="text-gray-400">Company ID</span>
                            <span className="font-mono text-xs text-indigo-600">{tenant.company_id}</span>
                        </div>
                        <div className="flex justify-between py-3">
                            <span className="text-gray-400">Node ID</span>
                            <span className="font-mono text-xs text-indigo-600">{tenant.id}</span>
                        </div>
                    </div>
                </div>

                <div className="bg-white rounded-3xl shadow-sm border border-gray-100 p-8">
                    <h3 className="text-xl font-bold text-gray-900 mb-6 flex items-center gap-2">
                        <Calendar size={20} className="text-indigo-600" />
                        Usage Analytics
                    </h3>
                    <div className="space-y-4 text-sm">
                        <div className="flex justify-between py-3 border-b border-gray-50">
                            <span className="text-gray-400">Total Trips Processed</span>
                            <span className="font-bold text-indigo-600">{tenant.trip_count}</span>
                        </div>
                        <div className="flex justify-between py-3 border-b border-gray-50">
                            <span className="text-gray-400">Active Sessions</span>
                            <span className="font-bold text-gray-700">12 (Current)</span>
                        </div>
                        <div className="flex justify-between py-3">
                            <span className="text-gray-400">Environment API Latency</span>
                            <span className="text-emerald-500 font-bold">42ms</span>
                        </div>
                    </div>
                </div>
            </div>

            <div className="mt-8 bg-indigo-900 rounded-3xl p-8 border border-white/10 relative overflow-hidden group">
                <div className="absolute top-0 right-0 p-8 text-white/5 pointer-events-none group-hover:scale-110 transition-transform">
                    <Shield size={120} />
                </div>
                <div className="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
                    <div>
                        <h3 className="text-2xl font-bold text-white mb-2 flex items-center gap-3">
                            <Key className="text-indigo-400" />
                            Identity & Access
                        </h3>
                        <p className="text-indigo-200 text-sm max-w-lg">
                            Provision administrative credentials for this node. New administrators will be required 
                            to reset their passwords upon first login to ensure cryptographic security.
                        </p>
                    </div>
                    <button 
                        onClick={() => setIsProvisionModalOpen(true)}
                        className="px-6 py-3 bg-white text-indigo-900 rounded-2xl font-bold hover:bg-indigo-50 transition-all shadow-xl shadow-black/20 flex items-center gap-2 active:scale-95"
                    >
                        Provision Admin Authority
                    </button>
                </div>
            </div>

            {/* Provision Modal */}
            {isProvisionModalOpen && (
                <div className="fixed inset-0 bg-indigo-950/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
                    <div className="bg-white rounded-[2.5rem] w-full max-w-md p-8 shadow-2xl relative animate-in fade-in zoom-in duration-200">
                        <button 
                            onClick={() => setIsProvisionModalOpen(false)}
                            className="absolute top-6 right-6 p-2 text-gray-400 hover:text-gray-600 rounded-full hover:bg-gray-100 transition-colors"
                        >
                            <X size={20} />
                        </button>

                        <div className="mb-8">
                            <div className="w-16 h-16 bg-indigo-50 text-indigo-600 rounded-2xl flex items-center justify-center mb-4">
                                <UsersIcon size={32} />
                            </div>
                            <h2 className="text-2xl font-black text-gray-900">Provision Admin</h2>
                            <p className="text-gray-500 text-sm">Deploy identity for {tenant.name}</p>
                        </div>

                        <form onSubmit={handleProvisionAdmin} className="space-y-6">
                            <div>
                                <label className="block text-xs font-black text-gray-400 uppercase tracking-widest mb-2 ml-1">Full Name</label>
                                <div className="relative">
                                    <UsersIcon className="absolute left-4 top-3.5 text-gray-400" size={18} />
                                    <input 
                                        type="text" 
                                        required
                                        value={adminName}
                                        onChange={(e) => setAdminName(e.target.value)}
                                        className="w-full bg-gray-50 border-0 rounded-2xl py-3.5 pl-12 pr-4 focus:ring-2 focus:ring-indigo-600 transition-all font-medium text-gray-900"
                                        placeholder="e.g. John Doe"
                                    />
                                </div>
                            </div>

                            <div>
                                <label className="block text-xs font-black text-gray-400 uppercase tracking-widest mb-2 ml-1">Corporate Email</label>
                                <div className="relative">
                                    <Mail className="absolute left-4 top-3.5 text-gray-400" size={18} />
                                    <input 
                                        type="email" 
                                        required
                                        value={adminEmail}
                                        onChange={(e) => setAdminEmail(e.target.value)}
                                        className="w-full bg-gray-50 border-0 rounded-2xl py-3.5 pl-12 pr-4 focus:ring-2 focus:ring-indigo-600 transition-all font-medium text-gray-900"
                                        placeholder="admin@aad-limousine.ae"
                                    />
                                </div>
                            </div>

                            <div>
                                <label className="block text-xs font-black text-gray-400 uppercase tracking-widest mb-2 ml-1">Initial Secret</label>
                                <div className="relative">
                                    <Key className="absolute left-4 top-3.5 text-gray-400" size={18} />
                                    <input 
                                        type="password" 
                                        required
                                        value={adminPassword}
                                        onChange={(e) => setAdminPassword(e.target.value)}
                                        className="w-full bg-gray-50 border-0 rounded-2xl py-3.5 pl-12 pr-4 focus:ring-2 focus:ring-indigo-600 transition-all font-medium text-gray-900"
                                        placeholder="••••••••••••"
                                    />
                                </div>
                            </div>

                            <button 
                                type="submit"
                                disabled={isSubmitting}
                                className="w-full bg-indigo-600 text-white rounded-2xl py-4 font-bold hover:bg-indigo-700 transition-all shadow-xl shadow-indigo-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                            >
                                {isSubmitting ? 'Deploying Access...' : 'Generate Admin Account'}
                            </button>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default TenantDetail;
