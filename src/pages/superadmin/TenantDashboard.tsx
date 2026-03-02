import { useState, useEffect } from 'react';
import { superAdminService, TenantWithFeatures } from '../../services/superadmin';
import { 
    Globe, 
    ShieldCheck, 
    Settings, 
    Search, 
    Database, 
    Activity, 
    Zap,
    Filter,
    Plus,
    X,
    Calendar,
    ArrowUpRight
} from 'lucide-react';

const FEATURE_LABELS: Record<string, string> = {
    'dispatch': 'Internal Dispatch',
    'outsource_marketplace': 'B2B Marketplace',
    'fleet_management': 'Fleet Control',
    'team_management': 'Team Manager',
    'invoicing': 'Smart Invoicing',
    'partner_portal': 'Partner Access',
    'analytics': 'BI Analytics'
};

const TenantDashboard = () => {
    const [tenants, setTenants] = useState<TenantWithFeatures[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState('');
    const [selectedTenant, setSelectedTenant] = useState<TenantWithFeatures | null>(null);
    const [isFeatureModalOpen, setFeatureModalOpen] = useState(false);
    const [isCreateModalOpen, setCreateModalOpen] = useState(false);
    const [newTenant, setNewTenant] = useState({ name: '', slug: '', plan: 'STARTER' });

    useEffect(() => {
        fetchTenants();
    }, []);

    const fetchTenants = async () => {
        try {
            setLoading(true);
            const data = await superAdminService.listTenants();
            setTenants(data);
        } catch (error) {
            console.error('Failed to fetch tenants:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleCreateTenant = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await superAdminService.createTenant(newTenant);
            setCreateModalOpen(false);
            setNewTenant({ name: '', slug: '', plan: 'STARTER' });
            fetchTenants();
        } catch (error) {
            console.error('Failed to create tenant:', error);
        }
    };

    const handleSwitchTenant = async (tenantId: string) => {
        try {
            await superAdminService.switchTenant(tenantId);
        } catch (error) {
            console.error('Failed to switch tenant:', error);
        }
    };

    const handleToggleFeature = async (tenantId: string, featureKey: string, currentEnabled: boolean) => {
        try {
            await superAdminService.toggleFeature(tenantId, featureKey, !currentEnabled);
            fetchTenants();
            if (selectedTenant && selectedTenant.id === tenantId) {
                const refreshed = await superAdminService.listTenants();
                const updated = refreshed.find(t => t.id === tenantId);
                if (updated) setSelectedTenant(updated);
            }
        } catch (error) {
            console.error('Failed to toggle feature:', error);
        }
    };

    const filteredTenants = tenants.filter(t => 
        t.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
        t.slug.toLowerCase().includes(searchQuery.toLowerCase())
    );

    const stats = {
        totalTenants: tenants.length,
        activeTenants: tenants.filter(t => t.status === 'ACTIVE').length,
        totalPremium: tenants.filter(t => t.plan === 'ENTERPRISE' || t.plan === 'PROFESSIONAL').length,
        totalTrips: tenants.reduce((acc, t) => acc + t.trip_count, 0)
    };

    if (loading && tenants.length === 0) {
        return (
            <div className="flex items-center justify-center min-h-[400px]">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
            </div>
        );
    }

    return (
        <div className="p-6 max-w-7xl mx-auto animate-in fade-in duration-500">
            {/* Header */}
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900 tracking-tight">System Orchestration</h1>
                    <p className="text-gray-500 mt-1">Multi-tenant management and granular microservice control</p>
                </div>
                <div className="flex items-center gap-3">
                    <button className="flex items-center gap-2 px-4 py-2 bg-white border border-gray-200 rounded-xl text-gray-700 font-medium hover:bg-gray-50 transition-all shadow-sm">
                        <Filter size={18} />
                        Filters
                    </button>
                    <button 
                        onClick={() => setCreateModalOpen(true)}
                        className="flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-xl font-medium hover:bg-indigo-700 transition-all shadow-md active:scale-95"
                    >
                        <Plus size={18} />
                        New Tenant
                    </button>
                </div>
            </div>

            {/* Stats Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                {[
                    { label: 'Cloud Tenants', value: stats.totalTenants, icon: Globe, color: 'indigo' },
                    { label: 'Operational Nodes', value: stats.activeTenants, icon: Activity, color: 'emerald' },
                    { label: 'System Transactions', value: stats.totalTrips, icon: Zap, color: 'amber' },
                    { label: 'Premium Nodes', value: stats.totalPremium, icon: ShieldCheck, color: 'violet' }
                ].map((stat, i) => (
                    <div key={i} className="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition-all group">
                        <div className="flex items-center justify-between mb-4">
                            <div className={`p-2.5 rounded-xl bg-${stat.color}-50 text-${stat.color}-600 group-hover:scale-110 transition-transform`}>
                                <stat.icon size={20} />
                            </div>
                            <span className="text-xs font-semibold text-gray-400 uppercase tracking-wider">{stat.label}</span>
                        </div>
                        <div className="flex items-baseline gap-2">
                            <h3 className="text-2xl font-bold text-gray-900">{stat.value}</h3>
                            <span className="text-xs text-emerald-500 font-medium flex items-center gap-0.5">
                                <ArrowUpRight size={12} />
                                +12%
                            </span>
                        </div>
                    </div>
                ))}
            </div>

            {/* Main Table Container */}
            <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
                <div className="p-5 border-b border-gray-50 flex flex-col md:flex-row md:items-center justify-between gap-4 bg-gray-50/30">
                    <div className="relative flex-1 max-w-md">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
                        <input 
                            type="text" 
                            placeholder="Identify tenant by name or node identifier..."
                            className="w-full pl-10 pr-4 py-2.5 bg-white border border-gray-200 rounded-xl focus:ring-2 focus:ring-indigo-500/10 focus:border-indigo-500 transition-all outline-none"
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                        />
                    </div>
                    <div className="flex items-center gap-4 text-sm text-gray-500">
                        <span className="flex items-center gap-1.5 font-medium tabular-nums">
                            <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></div>
                            {tenants.length} system nodes discovered
                        </span>
                    </div>
                </div>

                <div className="overflow-x-auto">
                    <table className="w-full text-left">
                        <thead>
                            <tr className="bg-gray-50 text-gray-500 text-xs font-bold uppercase tracking-wider">
                                <th className="px-6 py-4">Tenant Identity</th>
                                <th className="px-6 py-4">Lifecycle State</th>
                                <th className="px-6 py-4">Allocation Plan</th>
                                <th className="px-6 py-4">Capacity</th>
                                <th className="px-6 py-4">Microservices</th>
                                <th className="px-6 py-4 text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-50">
                            {filteredTenants.map((tenant) => (
                                <tr key={tenant.id} className="hover:bg-indigo-50/30 transition-colors group">
                                    <td className="px-6 py-4">
                                        <div className="flex items-center gap-3">
                                            <div className="w-10 h-10 rounded-xl bg-indigo-50 flex items-center justify-center text-indigo-600 font-bold border border-indigo-100 shadow-sm">
                                                {tenant.name.charAt(0)}
                                            </div>
                                            <div>
                                                <p className="font-semibold text-gray-900 group-hover:text-indigo-600 transition-colors">{tenant.name}</p>
                                                <p className="text-xs text-gray-400 font-mono tracking-tighter">{tenant.slug}</p>
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-6 py-4">
                                        <span className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold ${
                                            tenant.status === 'ACTIVE' 
                                                ? 'bg-emerald-50 text-emerald-700 border border-emerald-100' 
                                                : 'bg-amber-50 text-amber-700 border border-amber-100'
                                        }`}>
                                            <div className={`w-1.5 h-1.5 rounded-full ${tenant.status === 'ACTIVE' ? 'bg-emerald-500' : 'bg-amber-500'}`}></div>
                                            {tenant.status}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="flex flex-col">
                                            <span className="text-sm font-medium text-gray-700">{tenant.plan}</span>
                                            <span className="text-[10px] text-gray-400 flex items-center gap-1">
                                                <Calendar size={10} />
                                                Next Renewal: Mar 20, 2026
                                            </span>
                                        </div>
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="flex items-center gap-4">
                                            <div className="flex flex-col">
                                                <span className="text-xs font-bold text-gray-400 uppercase tracking-tighter">System Users</span>
                                                <span className="text-sm font-medium text-gray-800 tabular-nums">{tenant.user_count} / {tenant.max_users}</span>
                                            </div>
                                            <div className="flex flex-col">
                                                <span className="text-xs font-bold text-gray-400 uppercase tracking-tighter">Transactions</span>
                                                <span className="text-sm font-medium text-gray-800 tabular-nums">{tenant.trip_count}</span>
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="flex -space-x-1.5 overflow-hidden">
                                            {tenant.features.filter(f => f.is_enabled).slice(0, 4).map((f, i) => (
                                                <div key={i} className="inline-flex items-center justify-center w-7 h-7 rounded-lg bg-white border border-gray-100 text-indigo-600 shadow-sm ring-2 ring-white" title={FEATURE_LABELS[f.feature_key]}>
                                                    <Zap size={12} fill="currentColor" className="opacity-80" />
                                                </div>
                                            ))}
                                            {tenant.features.filter(f => f.is_enabled).length > 4 && (
                                                <div className="inline-flex items-center justify-center w-7 h-7 rounded-lg bg-gray-50 border border-gray-100 text-gray-400 font-bold text-[10px] ring-2 ring-white">
                                                    +{tenant.features.filter(f => f.is_enabled).length - 4}
                                                </div>
                                            )}
                                        </div>
                                    </td>
                                    <td className="px-6 py-4 text-right">
                                        <div className="flex items-center justify-end gap-2">
                                            <button 
                                                onClick={() => {
                                                    setSelectedTenant(tenant);
                                                    setFeatureModalOpen(true);
                                                }}
                                                className="p-2 text-gray-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-all"
                                                title="Environment Config"
                                            >
                                                <Settings size={18} />
                                            </button>
                                            <button 
                                                onClick={() => handleSwitchTenant(tenant.id)}
                                                className="flex items-center gap-1.5 px-3 py-1.5 bg-indigo-50 text-indigo-600 rounded-lg font-bold text-xs hover:bg-indigo-600 hover:text-white transition-all group/btn"
                                            >
                                                Enter
                                                <ArrowUpRight size={14} className="group-hover/btn:translate-x-0.5 group-hover/btn:-translate-y-0.5 transition-transform" />
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>

            {/* Create Tenant Modal */}
            {isCreateModalOpen && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/60 backdrop-blur-sm animate-in fade-in zoom-in duration-300">
                    <div className="bg-white rounded-3xl w-full max-w-md shadow-2xl overflow-hidden">
                        <div className="p-6 border-b border-gray-100 flex items-center justify-between">
                            <h2 className="text-xl font-bold text-gray-900">Provision New Tenant</h2>
                            <button onClick={() => setCreateModalOpen(false)} className="text-gray-400 hover:text-gray-600">
                                <X size={20} />
                            </button>
                        </div>
                        <form onSubmit={handleCreateTenant} className="p-6 space-y-4">
                            <div>
                                <label className="block text-xs font-bold text-gray-400 uppercase mb-1">Company Name</label>
                                <input 
                                    type="text"
                                    required
                                    className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:ring-2 focus:ring-indigo-500/10 focus:border-indigo-500 outline-none transition-all"
                                    value={newTenant.name}
                                    onChange={(e) => setNewTenant({...newTenant, name: e.target.value})}
                                />
                            </div>
                            <div>
                                <label className="block text-xs font-bold text-gray-400 uppercase mb-1">Node Identifier (Slug)</label>
                                <input 
                                    type="text"
                                    required
                                    placeholder="e.g. apex-limo"
                                    className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:ring-2 focus:ring-indigo-500/10 focus:border-indigo-500 outline-none transition-all"
                                    value={newTenant.slug}
                                    onChange={(e) => setNewTenant({...newTenant, slug: e.target.value.toLowerCase().replace(/\s+/g, '-')})}
                                />
                            </div>
                            <div>
                                <label className="block text-xs font-bold text-gray-400 uppercase mb-1">Allocation Plan</label>
                                <select 
                                    className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:ring-2 focus:ring-indigo-500/10 focus:border-indigo-500 outline-none transition-all"
                                    value={newTenant.plan}
                                    onChange={(e) => setNewTenant({...newTenant, plan: e.target.value})}
                                >
                                    <option value="STARTER">Starter Pack</option>
                                    <option value="PROFESSIONAL">Professional Tier</option>
                                    <option value="ENTERPRISE">Enterprise Node</option>
                                </select>
                            </div>
                            <button 
                                type="submit"
                                className="w-full py-3 bg-indigo-600 text-white rounded-xl font-bold shadow-lg hover:bg-indigo-700 transition-all active:scale-95 mt-4"
                            >
                                Provision Instance
                            </button>
                        </form>
                    </div>
                </div>
            )}

            {/* Feature Gating Modal (Tenant Config) */}
            {isFeatureModalOpen && selectedTenant && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/60 backdrop-blur-sm animate-in fade-in zoom-in duration-300">
                    <div className="bg-white rounded-3xl w-full max-w-2xl shadow-2xl overflow-hidden border border-white/20">
                        {/* Modal Header */}
                        <div className="p-6 bg-gradient-to-br from-indigo-600 to-indigo-800 text-white relative">
                            <button 
                                onClick={() => setFeatureModalOpen(false)}
                                className="absolute top-4 right-4 p-2 bg-white/10 hover:bg-white/20 rounded-xl transition-all"
                            >
                                <X size={20} />
                            </button>
                            <div className="flex items-center gap-4">
                                <div className="w-14 h-14 rounded-2xl bg-white/10 backdrop-blur-md flex items-center justify-center text-2xl font-bold shadow-inner">
                                    {selectedTenant.name.charAt(0)}
                                </div>
                                <div>
                                    <h2 className="text-xl font-bold tracking-tight">Node Configuration: {selectedTenant.name}</h2>
                                    <p className="text-indigo-100 text-sm mt-0.5 opacity-80">Orchestrate microservice availability for this tenant instance</p>
                                </div>
                            </div>
                        </div>

                        {/* Modal Body */}
                        <div className="p-8">
                            <div className="flex items-center gap-2 mb-6 p-4 bg-indigo-50 rounded-2xl text-indigo-700 border border-indigo-100 shadow-sm">
                                <Database size={20} />
                                <span className="text-sm font-semibold tracking-wide">AVAILABLE MICROSERVICES</span>
                            </div>

                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                {selectedTenant.features.map((feature) => (
                                    <div 
                                        key={feature.feature_key}
                                        className={`flex items-center justify-between p-4 rounded-2xl border transition-all ${
                                            feature.is_enabled 
                                                ? 'bg-white border-indigo-100 shadow-sm' 
                                                : 'bg-gray-50 border-gray-100 opacity-70'
                                        }`}
                                    >
                                        <div className="flex items-center gap-3">
                                            <div className={`p-2 rounded-xl ${feature.is_enabled ? 'bg-indigo-50 text-indigo-600' : 'bg-gray-100 text-gray-400'}`}>
                                                <Zap size={16} />
                                            </div>
                                            <div>
                                                <p className={`text-sm font-bold ${feature.is_enabled ? 'text-gray-900' : 'text-gray-400'}`}>
                                                    {FEATURE_LABELS[feature.feature_key] || feature.feature_key}
                                                </p>
                                                <p className="text-[10px] text-gray-400 font-mono tracking-tighter uppercase">{feature.feature_key}</p>
                                            </div>
                                        </div>
                                        
                                        <button 
                                            onClick={() => handleToggleFeature(selectedTenant.id, feature.feature_key, feature.is_enabled)}
                                            className={`relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none ${
                                                feature.is_enabled ? 'bg-indigo-600' : 'bg-gray-200'
                                            }`}
                                        >
                                            <span className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${
                                                feature.is_enabled ? 'translate-x-5' : 'translate-x-0'
                                            }`} />
                                        </button>
                                    </div>
                                ))}
                            </div>
                        </div>

                        {/* Modal Footer */}
                        <div className="p-6 bg-gray-50 border-t border-gray-100 flex items-center justify-between">
                            <div className="flex items-center gap-4 text-xs text-gray-400">
                                <span className="flex items-center gap-1 font-medium">
                                    <ShieldCheck size={14} className="text-indigo-400" />
                                    Security Isolation: Active
                                </span>
                                <span className="flex items-center gap-1 font-medium">
                                    <Globe size={14} className="text-emerald-400" />
                                    Edge Synchronized
                                </span>
                            </div>
                            <button 
                                onClick={() => setFeatureModalOpen(false)}
                                className="px-6 py-2.5 bg-gray-900 text-white rounded-xl font-bold text-sm shadow-xl hover:bg-black transition-all active:scale-95"
                            >
                                Apply Configuration
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default TenantDashboard;
