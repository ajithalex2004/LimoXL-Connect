import { useState, useEffect } from 'react';
import { operatorService, type BookingConfig, type NUIMaster } from '../../services/api';
import {
    Plus, Search, Edit2, Trash2, X, Settings,
    Clock, Tool, Target, ShieldCheck,
    ArrowUpDown, Boxes, Layers, Tag as TagIcon
} from 'lucide-react';

const BookingConfigs = () => {
    const [configs, setConfigs] = useState<BookingConfig[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');

    // Masters for dropdowns
    const [vehicleClasses, setVehicleClasses] = useState<NUIMaster[]>([]);
    const [vehicleGroups, setVehicleGroups] = useState<NUIMaster[]>([]);
    const [vehicleUsages, setVehicleUsages] = useState<NUIMaster[]>([]);

    // Modal State
    const [showModal, setShowModal] = useState(false);
    const [editingConfig, setEditingConfig] = useState<BookingConfig | null>(null);
    const [formData, setFormData] = useState<Partial<BookingConfig>>({
        name: '',
        booking_type: 'Regular',
        request_type: 'On Demand',
        priority: 'Normal',
        sort_order: 1000,
        vehicle_classes: [],
        vehicle_groups: [],
        vehicle_usages: [],
        pickup_buffer: 60,
        auto_dispatch_buffer: 30,
        pricing_source: 'Zone Based',
        approval_workflow_required: false,
        epod_required: false,
        is_active: true
    });

    useEffect(() => {
        loadData();
        loadMasters();
    }, []);

    const loadData = async () => {
        try {
            setLoading(true);
            const data = await operatorService.listBookingConfigs();
            setConfigs(data || []);
        } catch (error) {
            console.error("Failed to load configs", error);
        } finally {
            setLoading(false);
        }
    };

    const loadMasters = async () => {
        try {
            const [classes, groups, usages] = await Promise.all([
                operatorService.listMasters('CLASS'),
                operatorService.listMasters('TYPE'),
                operatorService.listMasters('USAGE')
            ]);
            setVehicleClasses(classes || []);
            setVehicleGroups(groups || []);
            setVehicleUsages(usages || []);
        } catch (error) {
            console.error("Failed to load masters", error);
        }
    };

    const handleOpenAdd = () => {
        setEditingConfig(null);
        setFormData({
            name: '',
            booking_type: 'Regular',
            request_type: 'On Demand',
            priority: 'Normal',
            sort_order: 1000,
            vehicle_classes: [],
            vehicle_groups: [],
            vehicle_usages: [],
            pickup_buffer: 60,
            auto_dispatch_buffer: 30,
            pricing_source: 'Zone Based',
            approval_workflow_required: false,
            epod_required: false,
            is_active: true
        });
        setShowModal(true);
    };

    const handleOpenEdit = (config: BookingConfig) => {
        setEditingConfig(config);
        setFormData({ ...config });
        setShowModal(true);
    };

    const handleDelete = async (id: string) => {
        if (!confirm('Are you sure you want to delete this configuration?')) return;
        try {
            await operatorService.deleteBookingConfig(id);
            loadData();
        } catch (error) {
            console.error("Failed to delete config", error);
            alert("Failed to delete category");
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            if (editingConfig) {
                await operatorService.updateBookingConfig(editingConfig.id, formData);
            } else {
                await operatorService.createBookingConfig(formData);
            }
            setShowModal(false);
            loadData();
        } catch (error) {
            console.error("Failed to save config", error);
            alert("Failed to save configuration. Check for duplicate names.");
        }
    };

    const toggleArrayItem = (field: 'vehicle_classes' | 'vehicle_groups' | 'vehicle_usages', value: string) => {
        const current = (formData[field] as string[]) || [];
        const updated = current.includes(value)
            ? current.filter(v => v !== value)
            : [...current, value];
        setFormData({ ...formData, [field]: updated });
    };

    const filteredConfigs = configs.filter(c =>
        c.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        c.booking_type.toLowerCase().includes(searchTerm.toLowerCase())
    );

    return (
        <div className="space-y-6">
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">Booking Configuration</h1>
                    <p className="text-gray-500 text-sm">Manage business rules and dispatch logic for different booking types</p>
                </div>
                <div className="flex items-center gap-3">
                    <div className="relative">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
                        <input
                            type="text"
                            placeholder="Search configs..."
                            className="pl-10 pr-4 py-2 border border-gray-200 rounded-lg text-sm w-64 focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500"
                            value={searchTerm}
                            onChange={e => setSearchTerm(e.target.value)}
                        />
                    </div>
                    <button
                        onClick={handleOpenAdd}
                        className="flex items-center gap-2 px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors text-sm font-semibold shadow-sm"
                    >
                        <Plus className="h-4 w-4" /> Create New
                    </button>
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {loading ? (
                    <div className="col-span-full py-20 flex justify-center">
                        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-600"></div>
                    </div>
                ) : filteredConfigs.length === 0 ? (
                    <div className="col-span-full py-20 text-center bg-white rounded-2xl border-2 border-dashed border-gray-200">
                        <Settings className="h-12 w-12 text-gray-300 mx-auto mb-4" />
                        <h3 className="text-lg font-bold text-gray-900">No configurations found</h3>
                        <p className="text-gray-500">Get started by creating your first booking configuration</p>
                    </div>
                ) : (
                    filteredConfigs.map((config) => (
                        <div key={config.id} className="bg-white rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition-shadow overflow-hidden group">
                            <div className="p-5">
                                <div className="flex items-start justify-between mb-4">
                                    <div className="flex items-center gap-3">
                                        <div className="p-2.5 bg-emerald-50 text-emerald-600 rounded-xl">
                                            <Tool className="h-5 w-5" />
                                        </div>
                                        <div>
                                            <h3 className="font-bold text-gray-900 leading-none mb-1">{config.name}</h3>
                                            <div className="flex items-center gap-2">
                                                <span className="text-[10px] font-bold text-gray-400 uppercase tracking-tighter">Sort Order: {config.sort_order}</span>
                                                <span className={`w-1.5 h-1.5 rounded-full ${config.is_active ? 'bg-emerald-500' : 'bg-gray-300'}`}></span>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                                        <button onClick={() => handleOpenEdit(config)} className="p-1.5 text-gray-400 hover:text-emerald-600 hover:bg-emerald-50 rounded-lg">
                                            <Edit2 className="h-4 w-4" />
                                        </button>
                                        <button onClick={() => handleDelete(config.id)} className="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg">
                                            <Trash2 className="h-4 w-4" />
                                        </button>
                                    </div>
                                </div>

                                <div className="grid grid-cols-3 gap-2 mb-4">
                                    <div className="p-2 bg-gray-50 rounded-lg border border-gray-100">
                                        <div className="text-[10px] text-gray-400 font-bold uppercase mb-0.5">Type</div>
                                        <div className="text-xs font-bold text-gray-700 truncate">{config.booking_type}</div>
                                    </div>
                                    <div className="p-2 bg-gray-50 rounded-lg border border-gray-100">
                                        <div className="text-[10px] text-gray-400 font-bold uppercase mb-0.5">Priority</div>
                                        <div className="text-xs font-bold text-gray-700 truncate">{config.priority}</div>
                                    </div>
                                    <div className="p-2 bg-gray-50 rounded-lg border border-gray-100">
                                        <div className="text-[10px] text-gray-400 font-bold uppercase mb-0.5">Price</div>
                                        <div className="text-xs font-bold text-gray-700 truncate">{config.pricing_source}</div>
                                    </div>
                                </div>

                                <div className="space-y-3">
                                    <div className="flex flex-wrap gap-1.5">
                                        {config.vehicle_classes?.map(c => (
                                            <span key={c} className="px-2 py-0.5 bg-blue-50 text-blue-700 text-[10px] font-bold rounded-md">
                                                {c}
                                            </span>
                                        ))}
                                        {config.vehicle_groups?.map(g => (
                                            <span key={g} className="px-2 py-0.5 bg-purple-50 text-purple-700 text-[10px] font-bold rounded-md">
                                                {g}
                                            </span>
                                        ))}
                                    </div>

                                    <div className="flex items-center justify-between text-[11px] font-bold">
                                        <div className="flex items-center gap-1 text-gray-500">
                                            <Clock className="h-3 w-3" /> Pickup: {config.pickup_buffer}m
                                        </div>
                                        <div className="flex items-center gap-1 text-gray-500">
                                            <Target className="h-3 w-3" /> Dispatch: {config.auto_dispatch_buffer}m
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    ))
                )}
            </div>

            {/* Modal */}
            {showModal && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm animate-in fade-in duration-200">
                    <div className="bg-white rounded-3xl shadow-2xl max-w-4xl w-full flex flex-col overflow-hidden max-h-[90vh]">
                        {/* Modal Header */}
                        <div className="px-8 py-6 border-b border-gray-100 flex items-center justify-between bg-white sticky top-0 z-10">
                            <div>
                                <h2 className="text-xl font-bold text-gray-900">{editingConfig ? 'Edit Booking Configuration' : 'Create New Booking Configuration'}</h2>
                                <p className="text-xs text-gray-500 mt-1">Configure business rules for this booking type</p>
                            </div>
                            <button onClick={() => setShowModal(false)} className="p-2 hover:bg-gray-100 rounded-xl transition-colors">
                                <X className="h-5 w-5 text-gray-500" />
                            </button>
                        </div>

                        {/* Modal Body */}
                        <div className="flex-1 overflow-y-auto p-8 bg-gray-50/30">
                            <form id="configForm" onSubmit={handleSubmit} className="space-y-8">
                                {/* Basic Configuration */}
                                <div className="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm space-y-6">
                                    <div className="flex items-center gap-2 pb-2 border-b border-gray-50">
                                        <Settings className="h-4 w-4 text-emerald-600" />
                                        <h3 className="text-sm font-bold text-gray-900">Basic Configuration</h3>
                                    </div>
                                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                                        <div className="space-y-1.5">
                                            <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Configuration Name *</label>
                                            <input
                                                type="text"
                                                required
                                                placeholder="e.g., VIP Airport Transfer"
                                                className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                value={formData.name}
                                                onChange={e => setFormData({ ...formData, name: e.target.value })}
                                            />
                                        </div>
                                        <div className="space-y-1.5">
                                            <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Booking Type</label>
                                            <select
                                                className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm appearance-none"
                                                value={formData.booking_type}
                                                onChange={e => setFormData({ ...formData, booking_type: e.target.value })}
                                            >
                                                <option value="Regular">Regular</option>
                                                <option value="VIP">VIP</option>
                                                <option value="Staff">Staff</option>
                                                <option value="Partner">Partner</option>
                                            </select>
                                        </div>
                                        <div className="space-y-1.5">
                                            <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Priority</label>
                                            <select
                                                className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm appearance-none"
                                                value={formData.priority}
                                                onChange={e => setFormData({ ...formData, priority: e.target.value })}
                                            >
                                                <option value="Low">Low</option>
                                                <option value="Normal">Normal</option>
                                                <option value="High">High</option>
                                                <option value="Urgent">Urgent</option>
                                            </select>
                                        </div>
                                        <div className="space-y-1.5">
                                            <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Sort Order</label>
                                            <input
                                                type="number"
                                                className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                value={formData.sort_order}
                                                onChange={e => setFormData({ ...formData, sort_order: parseInt(e.target.value) })}
                                            />
                                            <p className="text-[10px] text-gray-400 mt-1">Order in which this configuration should appear</p>
                                        </div>
                                        <div className="space-y-1.5">
                                            <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Request Type</label>
                                            <select
                                                className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm appearance-none"
                                                value={formData.request_type}
                                                onChange={e => setFormData({ ...formData, request_type: e.target.value })}
                                            >
                                                <option value="On Demand">On Demand</option>
                                                <option value="Scheduled">Scheduled</option>
                                                <option value="Hourly">Hourly</option>
                                            </select>
                                        </div>
                                    </div>
                                </div>

                                {/* Vehicle Configuration */}
                                <div className="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm space-y-6">
                                    <div className="flex items-center gap-2 pb-2 border-b border-gray-50">
                                        <Boxes className="h-4 w-4 text-emerald-600" />
                                        <h3 className="text-sm font-bold text-gray-900">Vehicle Configuration</h3>
                                    </div>
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                                        <div className="space-y-3">
                                            <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider flex items-center gap-1.5">
                                                <Layers className="h-3 w-3" /> Vehicle Classes *
                                            </label>
                                            <div className="flex flex-wrap gap-2">
                                                {vehicleClasses.map(m => (
                                                    <button
                                                        key={m.id}
                                                        type="button"
                                                        onClick={() => toggleArrayItem('vehicle_classes', m.name)}
                                                        className={`px-3 py-1.5 rounded-lg text-xs font-bold border transition-all ${formData.vehicle_classes?.includes(m.name)
                                                            ? 'bg-emerald-600 text-white border-emerald-600 shadow-sm shadow-emerald-200'
                                                            : 'bg-white text-gray-600 border-gray-200 hover:border-emerald-200'
                                                            }`}
                                                    >
                                                        {m.name}
                                                    </button>
                                                ))}
                                            </div>
                                        </div>
                                        <div className="space-y-3">
                                            <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider flex items-center gap-1.5">
                                                <Tool className="h-3 w-3" /> Vehicle Types *
                                            </label>
                                            <div className="flex flex-wrap gap-2">
                                                {vehicleGroups.map(m => (
                                                    <button
                                                        key={m.id}
                                                        type="button"
                                                        onClick={() => toggleArrayItem('vehicle_groups', m.name)}
                                                        className={`px-3 py-1.5 rounded-lg text-xs font-bold border transition-all ${formData.vehicle_groups?.includes(m.name)
                                                            ? 'bg-emerald-600 text-white border-emerald-600 shadow-sm shadow-emerald-200'
                                                            : 'bg-white text-gray-600 border-gray-200 hover:border-emerald-200'
                                                            }`}
                                                    >
                                                        {m.name}
                                                    </button>
                                                ))}
                                            </div>
                                        </div>
                                        <div className="space-y-3 md:col-span-2">
                                            <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider flex items-center gap-1.5">
                                                <TagIcon className="h-3 w-3" /> Vehicle Usages *
                                            </label>
                                            <div className="flex flex-wrap gap-2">
                                                {vehicleUsages.map(m => (
                                                    <button
                                                        key={m.id}
                                                        type="button"
                                                        onClick={() => toggleArrayItem('vehicle_usages', m.name)}
                                                        className={`px-3 py-1.5 rounded-lg text-xs font-bold border transition-all ${formData.vehicle_usages?.includes(m.name)
                                                            ? 'bg-emerald-600 text-white border-emerald-600 shadow-sm shadow-emerald-200'
                                                            : 'bg-white text-gray-600 border-gray-200 hover:border-emerald-200'
                                                            }`}
                                                    >
                                                        {m.name}
                                                    </button>
                                                ))}
                                            </div>
                                        </div>
                                    </div>
                                </div>

                                {/* Business Rules */}
                                <div className="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm space-y-6">
                                    <div className="flex items-center gap-2 pb-2 border-b border-gray-50">
                                        <ShieldCheck className="h-4 w-4 text-emerald-600" />
                                        <h3 className="text-sm font-bold text-gray-900">Business Rules</h3>
                                    </div>
                                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                                        <div className="space-y-1.5">
                                            <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Pickup Buffer (Mins)</label>
                                            <input
                                                type="number"
                                                className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                value={formData.pickup_buffer}
                                                onChange={e => setFormData({ ...formData, pickup_buffer: parseInt(e.target.value) })}
                                            />
                                            <p className="text-[10px] text-gray-400 mt-1">Minimum time between booking and pickup</p>
                                        </div>
                                        <div className="space-y-1.5">
                                            <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Auto Dispatch Buffer (Mins)</label>
                                            <input
                                                type="number"
                                                className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                value={formData.auto_dispatch_buffer}
                                                onChange={e => setFormData({ ...formData, auto_dispatch_buffer: parseInt(e.target.value) })}
                                            />
                                            <p className="text-[10px] text-gray-400 mt-1">Buffer time for automatic driver dispatch</p>
                                        </div>
                                        <div className="space-y-1.5">
                                            <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Pricing Source</label>
                                            <select
                                                className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm appearance-none"
                                                value={formData.pricing_source}
                                                onChange={e => setFormData({ ...formData, pricing_source: e.target.value })}
                                            >
                                                <option value="Zone Based">Zone Based</option>
                                                <option value="Distance Based">Distance Based</option>
                                                <option value="Fixed Rate">Fixed Rate</option>
                                            </select>
                                        </div>
                                    </div>
                                    <div className="flex flex-wrap gap-8 pt-4">
                                        <label className="flex items-center gap-3 cursor-pointer group">
                                            <div className="relative">
                                                <input
                                                    type="checkbox"
                                                    className="sr-only"
                                                    checked={formData.approval_workflow_required}
                                                    onChange={e => setFormData({ ...formData, approval_workflow_required: e.target.checked })}
                                                />
                                                <div className={`w-10 h-5 rounded-full transition-colors ${formData.approval_workflow_required ? 'bg-emerald-600' : 'bg-gray-200'}`}></div>
                                                <div className={`absolute top-1 left-1 w-3 h-3 bg-white rounded-full transition-transform ${formData.approval_workflow_required ? 'translate-x-5' : ''}`}></div>
                                            </div>
                                            <span className="text-sm font-semibold text-gray-700">Approval Workflow Required</span>
                                        </label>
                                        <label className="flex items-center gap-3 cursor-pointer group">
                                            <div className="relative">
                                                <input
                                                    type="checkbox"
                                                    className="sr-only"
                                                    checked={formData.epod_required}
                                                    onChange={e => setFormData({ ...formData, epod_required: e.target.checked })}
                                                />
                                                <div className={`w-10 h-5 rounded-full transition-colors ${formData.epod_required ? 'bg-emerald-600' : 'bg-gray-200'}`}></div>
                                                <div className={`absolute top-1 left-1 w-3 h-3 bg-white rounded-full transition-transform ${formData.epod_required ? 'translate-x-5' : ''}`}></div>
                                            </div>
                                            <span className="text-sm font-semibold text-gray-700">EPOD Required</span>
                                        </label>
                                    </div>
                                </div>
                            </form>
                        </div>

                        {/* Modal Footer */}
                        <div className="px-8 py-6 border-t border-gray-100 bg-white flex items-center justify-end gap-3 sticky bottom-0 z-10">
                            <button
                                type="button"
                                onClick={() => setShowModal(false)}
                                className="px-6 py-2.5 text-sm font-bold text-gray-700 hover:bg-gray-50 border border-transparent rounded-xl transition-all"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                form="configForm"
                                className="px-10 py-2.5 bg-emerald-600 text-white text-sm font-bold rounded-xl hover:bg-emerald-700 transition-all shadow-lg shadow-emerald-600/20 active:scale-95"
                            >
                                {editingConfig ? 'Update Configuration' : 'Create Booking Configuration'}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default BookingConfigs;
