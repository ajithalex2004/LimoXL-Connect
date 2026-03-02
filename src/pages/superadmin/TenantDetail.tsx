import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { superAdminService, TenantWithFeatures } from '../../services/superadmin';
import { ArrowLeft, Building2, Calendar, Globe, Shield, Users as UsersIcon, Zap } from 'lucide-react';

const TenantDetail = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [tenant, setTenant] = useState<TenantWithFeatures | null>(null);
    const [loading, setLoading] = useState(true);

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
        </div>
    );
};

export default TenantDetail;
