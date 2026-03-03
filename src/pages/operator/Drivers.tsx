import { useState, useEffect } from 'react';
import { operatorService, type Driver, type FleetAttachment } from '../../services/api';
import { 
    Plus, User, Search, Edit2, Trash2, 
    X, Info, Paperclip, CheckCircle2,
    Calendar, Shield, Phone, Globe, Hash,
    Download, FileText
} from 'lucide-react';

const Drivers = () => {
    // List State
    const [drivers, setDrivers] = useState<Driver[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');
    
    // Modal State
    const [showModal, setShowModal] = useState(false);
    const [activeTab, setActiveTab] = useState<'personal' | 'attachments'>('personal');
    const [editingDriver, setEditingDriver] = useState<Driver | null>(null);
    const [attachments, setAttachments] = useState<FleetAttachment[]>([]);
    const [attLoading, setAttLoading] = useState(false);

    // Form State
    const [formData, setFormData] = useState<Partial<Driver>>({
        name: '',
        first_name: '',
        last_name: '',
        phone: '',
        license_number: '',
        license_expiry: '',
        itc_permit_expiry: '',
        visa_expiry: '',
        date_of_birth: '',
        nationality: '',
        emirates_id: '',
        date_of_join: '',
        dallas_id: '',
        communication_language: 'English',
        hierarchy: '',
        driver_type: 'Regular'
    });

    useEffect(() => {
        loadData();
    }, []);

    // Auto-generate name and full name
    useEffect(() => {
        const fullName = `${formData.first_name || ''} ${formData.last_name || ''}`.trim();
        if (formData.name !== fullName) {
            setFormData(prev => ({ ...prev, name: fullName }));
        }
    }, [formData.first_name, formData.last_name]);

    const loadData = async () => {
        try {
            setLoading(true);
            const data = await operatorService.listDrivers();
            setDrivers(data || []);
        } catch (error) {
            console.error("Failed to load drivers", error);
        } finally {
            setLoading(false);
        }
    };

    const handleOpenAdd = () => {
        setEditingDriver(null);
        setFormData({
            name: '',
            first_name: '',
            last_name: '',
            phone: '',
            license_number: '',
            license_expiry: '',
            itc_permit_expiry: '',
            visa_expiry: '',
            date_of_birth: '',
            nationality: '',
            emirates_id: '',
            date_of_join: '',
            dallas_id: '',
            communication_language: 'English',
            hierarchy: '',
            driver_type: 'Regular'
        });
        setActiveTab('personal');
        setShowModal(true);
    };

    const handleOpenEdit = async (driver: Driver) => {
        setEditingDriver(driver);
        setFormData({ ...driver });
        setActiveTab('personal');
        setShowModal(true);

        // Load attachments
        try {
            setAttLoading(true);
            const atts = await operatorService.listAttachments(driver.id, 'driver');
            setAttachments(atts || []);
        } catch (error) {
            console.error("Failed to load attachments", error);
        } finally {
            setAttLoading(false);
        }
    };

    const handleDelete = async (id: string) => {
        if (!confirm('Are you sure you want to delete this driver?')) return;
        try {
            await operatorService.deleteDriver(id);
            loadData();
        } catch (error) {
            console.error("Failed to delete driver", error);
            alert("Failed to delete driver");
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            if (editingDriver) {
                await operatorService.updateDriver(editingDriver.id, formData);
            } else {
                await operatorService.createDriver(formData);
            }
            setShowModal(false);
            loadData();
        } catch (error) {
            console.error("Failed to save driver", error);
            alert("Failed to save driver");
        }
    };

    const filteredDrivers = drivers.filter(d => 
        d.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        d.phone.toLowerCase().includes(searchTerm.toLowerCase())
    );

    if (loading && drivers.length === 0) {
        return (
            <div className="flex items-center justify-center h-full">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-600"></div>
            </div>
        );
    }

    return (
        <div className="space-y-6">
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">Driver Master</h1>
                    <p className="text-gray-500 text-sm">Manage your driver registry, employment and compliance documents</p>
                </div>
                <div className="flex items-center gap-3">
                    <div className="relative">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
                        <input
                            type="text"
                            placeholder="Search by name or phone..."
                            className="pl-10 pr-4 py-2 border border-gray-200 rounded-lg text-sm w-64 focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500"
                            value={searchTerm}
                            onChange={e => setSearchTerm(e.target.value)}
                        />
                    </div>
                    <button
                        onClick={handleOpenAdd}
                        className="flex items-center gap-2 px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors text-sm font-semibold shadow-sm"
                    >
                        <Plus className="h-4 w-4" /> Add Driver
                    </button>
                </div>
            </div>

            {/* Driver List */}
            <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full text-left">
                        <thead className="bg-gray-50 border-b border-gray-200">
                            <tr>
                                <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase">Driver Info</th>
                                <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase">License/ID</th>
                                <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase">Employment</th>
                                <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase">Compliance</th>
                                <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100">
                            {filteredDrivers.length === 0 ? (
                                <tr>
                                    <td colSpan={5} className="px-6 py-12 text-center text-gray-500">
                                        No drivers found. Add your first driver to get started.
                                    </td>
                                </tr>
                            ) : (
                                filteredDrivers.map((driver) => (
                                    <tr key={driver.id} className="hover:bg-gray-50/50 transition-colors group">
                                        <td className="px-6 py-4">
                                            <div className="flex items-center gap-3">
                                                <div className="h-10 w-10 bg-emerald-50 rounded-full flex items-center justify-center text-emerald-600">
                                                    <User className="h-5 w-5" />
                                                </div>
                                                <div>
                                                    <div className="font-bold text-gray-900">{driver.name}</div>
                                                    <div className="text-xs text-gray-500">{driver.phone || 'No phone'}</div>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className="text-sm font-medium text-gray-700">{driver.license_number || '-'}</div>
                                            <div className="text-xs text-gray-400">ID: {driver.emirates_id || '-'}</div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className="text-sm text-gray-700">{driver.driver_type}</div>
                                            <div className="text-xs text-gray-400">{driver.hierarchy}</div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className="space-y-1">
                                                <div className="flex items-center gap-2">
                                                    <div className={`h-1.5 w-1.5 rounded-full ${new Date(driver.license_expiry || '') > new Date() ? 'bg-emerald-500' : 'bg-red-500'}`}></div>
                                                    <span className="text-[10px] text-gray-500 uppercase">License: {driver.license_expiry ? new Date(driver.license_expiry).toLocaleDateString() : 'N/A'}</span>
                                                </div>
                                                <div className="flex items-center gap-2">
                                                    <div className={`h-1.5 w-1.5 rounded-full ${new Date(driver.visa_expiry || '') > new Date() ? 'bg-emerald-500' : 'bg-red-500'}`}></div>
                                                    <span className="text-[10px] text-gray-500 uppercase">Visa: {driver.visa_expiry ? new Date(driver.visa_expiry).toLocaleDateString() : 'N/A'}</span>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4 text-right">
                                            <div className="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                                <button onClick={() => handleOpenEdit(driver)} className="p-1.5 text-gray-400 hover:text-emerald-600 hover:bg-emerald-50 rounded-lg">
                                                    <Edit2 className="h-4 w-4" />
                                                </button>
                                                <button onClick={() => handleDelete(driver.id)} className="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg">
                                                    <Trash2 className="h-4 w-4" />
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
            </div>

            {/* Modal */}
            {showModal && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm animate-in fade-in duration-200">
                    <div className="bg-white rounded-2xl shadow-2xl max-w-4xl w-full flex flex-col max-h-[90vh] overflow-hidden">
                        {/* Modal Header */}
                        <div className="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
                            <h2 className="text-xl font-bold text-gray-900">{editingDriver ? 'Edit Driver' : 'Create New Driver'}</h2>
                            <button onClick={() => setShowModal(false)} className="p-2 hover:bg-gray-100 rounded-lg transition-colors">
                                <X className="h-5 w-5 text-gray-500" />
                            </button>
                        </div>

                        {/* Tabs */}
                        <div className="px-6 flex border-b border-gray-100">
                            <button
                                onClick={() => setActiveTab('personal')}
                                className={`px-4 py-3 text-sm font-semibold flex items-center gap-2 border-b-2 transition-colors ${
                                    activeTab === 'personal' ? 'border-emerald-600 text-emerald-600' : 'border-transparent text-gray-500 hover:text-gray-700'
                                }`}
                            >
                                <User className="h-4 w-4" /> Personal Details
                            </button>
                            <button
                                onClick={() => setActiveTab('attachments')}
                                className={`px-4 py-3 text-sm font-semibold flex items-center gap-2 border-b-2 transition-colors ${
                                    activeTab === 'attachments' ? 'border-emerald-600 text-emerald-600' : 'border-transparent text-gray-500 hover:text-gray-700'
                                }`}
                            >
                                <Paperclip className="h-4 w-4" /> Attachments
                            </button>
                        </div>

                        {/* Modal Body */}
                        <div className="flex-1 overflow-y-auto p-6">
                            <form id="driverForm" onSubmit={handleSubmit} className="space-y-8">
                                {activeTab === 'personal' ? (
                                    <>
                                        {/* Identity Section */}
                                        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                                            <div className="space-y-1.5">
                                                <label className="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                                                    <Hash className="h-3 w-3" /> Driver ID
                                                </label>
                                                <input
                                                    type="text"
                                                    placeholder="e.g., DRV001"
                                                    className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                    value={editingDriver?.id?.substring(0, 8) || 'AUTO-ID'}
                                                    readOnly
                                                />
                                            </div>
                                            <div className="space-y-1.5">
                                                <label className="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                                                    Hierarchy
                                                </label>
                                                <select
                                                    className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm appearance-none"
                                                    value={formData.hierarchy}
                                                    onChange={e => setFormData({ ...formData, hierarchy: e.target.value })}
                                                >
                                                    <option value="">Select hierarchy</option>
                                                    <option value="Executive">Executive</option>
                                                    <option value="Premium">Premium</option>
                                                </select>
                                            </div>
                                            <div className="space-y-1.5">
                                                <label className="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                                                    Driver Type
                                                </label>
                                                <select
                                                    className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm appearance-none"
                                                    value={formData.driver_type}
                                                    onChange={e => setFormData({ ...formData, driver_type: e.target.value })}
                                                >
                                                    <option value="Regular">Regular</option>
                                                    <option value="Temporary">Temporary</option>
                                                    <option value="VIP">VIP</option>
                                                </select>
                                            </div>
                                        </div>

                                        <div className="pt-4 space-y-4">
                                            <div className="flex items-center justify-between pb-2 border-b border-gray-100">
                                                <div className="flex items-center gap-2">
                                                    <User className="h-4 w-4 text-emerald-600" />
                                                    <span className="text-sm font-bold text-gray-900">Name Details</span>
                                                </div>
                                                <span className="text-[10px] text-amber-600 font-bold bg-amber-50 px-2 py-0.5 rounded">Required *</span>
                                            </div>
                                            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                                                <div className="space-y-1.5">
                                                    <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">First Name</label>
                                                    <input
                                                        type="text"
                                                        required
                                                        placeholder="First name (English)"
                                                        className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                        value={formData.first_name}
                                                        onChange={e => setFormData({ ...formData, first_name: e.target.value })}
                                                    />
                                                </div>
                                                <div className="space-y-1.5">
                                                    <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Last Name</label>
                                                    <input
                                                        type="text"
                                                        placeholder="Last name (English)"
                                                        className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                        value={formData.last_name}
                                                        onChange={e => setFormData({ ...formData, last_name: e.target.value })}
                                                    />
                                                </div>
                                                <div className="space-y-1.5">
                                                    <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Full Name</label>
                                                    <input
                                                        type="text"
                                                        readOnly
                                                        disabled
                                                        className="w-full px-4 py-2.5 bg-gray-100 border border-gray-200 rounded-xl text-gray-500 text-sm cursor-not-allowed"
                                                        value={formData.name}
                                                    />
                                                </div>
                                            </div>
                                        </div>

                                        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 pt-4">
                                            <div className="space-y-1.5">
                                                <label className="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                                                    <Calendar className="h-3 w-3" /> Date of Birth
                                                </label>
                                                <input
                                                    type="date"
                                                    className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                    value={formData.date_of_birth ? formData.date_of_birth.split('T')[0] : ''}
                                                    onChange={e => setFormData({ ...formData, date_of_birth: e.target.value })}
                                                />
                                            </div>
                                            <div className="space-y-1.5">
                                                <label className="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                                                    <Globe className="h-3 w-3" /> Nationality
                                                </label>
                                                <select
                                                    className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm appearance-none"
                                                    value={formData.nationality}
                                                    onChange={e => setFormData({ ...formData, nationality: e.target.value })}
                                                >
                                                    <option value="">Select nationality</option>
                                                    <option value="India">India</option>
                                                    <option value="Pakistan">Pakistan</option>
                                                    <option value="Egypt">Egypt</option>
                                                    <option value="UAE">UAE</option>
                                                </select>
                                            </div>
                                            <div className="space-y-1.5">
                                                <label className="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                                                    <Shield className="h-3 w-3" /> Emirates ID
                                                </label>
                                                <input
                                                    type="text"
                                                    placeholder="784-XXXX-XXXXXXX-X"
                                                    className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                    value={formData.emirates_id}
                                                    onChange={e => setFormData({ ...formData, emirates_id: e.target.value })}
                                                />
                                            </div>
                                        </div>

                                        {/* Employment Details */}
                                        <div className="pt-4 space-y-4">
                                            <div className="flex items-center gap-2 pb-2 border-b border-gray-100">
                                                <Calendar className="h-4 w-4 text-emerald-600" />
                                                <span className="text-sm font-bold text-gray-900">Employment Details</span>
                                            </div>
                                            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                                <div className="space-y-1.5">
                                                    <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Date of Join</label>
                                                    <input
                                                        type="date"
                                                        className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                        value={formData.date_of_join ? formData.date_of_join.split('T')[0] : ''}
                                                        onChange={e => setFormData({ ...formData, date_of_join: e.target.value })}
                                                    />
                                                </div>
                                                <div className="space-y-1.5">
                                                    <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Dallas ID</label>
                                                    <input
                                                        type="text"
                                                        placeholder="Enter Dallas ID"
                                                        className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                        value={formData.dallas_id}
                                                        onChange={e => setFormData({ ...formData, dallas_id: e.target.value })}
                                                    />
                                                </div>
                                            </div>
                                        </div>

                                        {/* Communication Section */}
                                        <div className="pt-4 space-y-4">
                                            <div className="flex items-center gap-2 pb-2 border-b border-gray-100">
                                                <Phone className="h-4 w-4 text-emerald-600" />
                                                <span className="text-sm font-bold text-gray-900">Contact & Communication</span>
                                            </div>
                                            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                                <div className="space-y-1.5">
                                                    <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Mobile Number *</label>
                                                    <div className="flex gap-2">
                                                        <div className="w-24 px-3 py-2.5 bg-gray-100 border border-gray-200 rounded-xl text-sm flex items-center gap-2">
                                                            <span className="text-lg">🇦er</span> +971
                                                        </div>
                                                        <input
                                                            type="tel"
                                                            required
                                                            placeholder="50 123 4567"
                                                            className="flex-1 px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                            value={formData.phone}
                                                            onChange={e => setFormData({ ...formData, phone: e.target.value })}
                                                        />
                                                    </div>
                                                </div>
                                                <div className="space-y-1.5">
                                                    <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Communication Language</label>
                                                    <select
                                                        className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm appearance-none"
                                                        value={formData.communication_language}
                                                        onChange={e => setFormData({ ...formData, communication_language: e.target.value })}
                                                    >
                                                        <option value="English">English</option>
                                                        <option value="Arabic">Arabic</option>
                                                        <option value="Hindi">Hindi</option>
                                                        <option value="Urdu">Urdu</option>
                                                    </select>
                                                </div>
                                            </div>
                                        </div>
                                    </>
                                ) : (
                                    <div className="space-y-6 animate-in slide-in-from-right-2 duration-300">
                                        {!editingDriver ? (
                                            <div className="p-12 text-center text-gray-500 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200">
                                                <Info className="mx-auto h-8 w-8 mb-2 opacity-20" />
                                                <p>Please save the driver first before adding attachments.</p>
                                            </div>
                                        ) : (
                                            <>
                                                <div className="p-12 border-2 border-dashed border-gray-200 rounded-2xl flex flex-col items-center justify-center gap-4 bg-gray-50/50">
                                                    <div className="p-4 bg-white rounded-2xl shadow-sm border border-gray-100">
                                                        <Paperclip className="h-8 w-8 text-gray-400" />
                                                    </div>
                                                    <div className="text-center">
                                                        <p className="font-bold text-gray-900">Upload Attachments</p>
                                                        <p className="text-xs text-gray-500 mt-1">License, Visa Copy, Emirates ID, Health Certificate etc.</p>
                                                    </div>
                                                    <button type="button" className="px-6 py-2 bg-white border border-gray-200 rounded-xl text-sm font-bold hover:bg-gray-50 transition-colors shadow-sm">
                                                        Select Files
                                                    </button>
                                                </div>
                                                
                                                <div className="space-y-3">
                                                    <p className="text-xs font-bold text-gray-500 uppercase tracking-wider">Attached Files ({attachments.length})</p>
                                                    {attLoading ? (
                                                        <div className="text-center py-4 text-gray-400 text-sm italic">Loading attachments...</div>
                                                    ) : attachments.length === 0 ? (
                                                        <div className="text-center py-8 text-gray-400 text-sm">No attachments yet.</div>
                                                    ) : (
                                                        <div className="grid grid-cols-1 gap-2">
                                                            {attachments.map((att) => (
                                                                <div key={att.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-xl border border-gray-100 group">
                                                                    <div className="flex items-center gap-3">
                                                                        <div className="h-8 w-8 bg-white rounded-lg flex items-center justify-center text-gray-400 shadow-sm">
                                                                            <FileText className="h-4 w-4" />
                                                                        </div>
                                                                        <div>
                                                                            <div className="text-sm font-semibold text-gray-700">{att.file_name}</div>
                                                                            <div className="text-[10px] text-gray-400 uppercase">{att.file_type} • {new Date(att.created_at || '').toLocaleDateString()}</div>
                                                                        </div>
                                                                    </div>
                                                                    <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                                                                        <a href={att.file_url} target="_blank" rel="noreferrer" className="p-1.5 text-gray-400 hover:text-emerald-600 rounded-lg">
                                                                            <Download className="h-4 w-4" />
                                                                        </a>
                                                                        <button 
                                                                            type="button" 
                                                                            onClick={async () => {
                                                                                if(confirm('Delete attachment?')) {
                                                                                    await operatorService.deleteAttachment(att.id);
                                                                                    const updated = await operatorService.listAttachments(editingDriver.id, 'driver');
                                                                                    setAttachments(updated);
                                                                                }
                                                                            }}
                                                                            className="p-1.5 text-gray-400 hover:text-red-600 rounded-lg"
                                                                        >
                                                                            <Trash2 className="h-4 w-4" />
                                                                        </button>
                                                                    </div>
                                                                </div>
                                                            ))}
                                                        </div>
                                                    )}
                                                </div>
                                            </>
                                        )}
                                    </div>
                                )}
                            </form>
                        </div>

                        {/* Modal Footer */}
                        <div className="px-6 py-4 border-t border-gray-100 flex items-center justify-end gap-3 bg-gray-50/30">
                            <button
                                type="button"
                                onClick={() => setShowModal(false)}
                                className="px-6 py-2.5 text-sm font-bold text-gray-700 hover:bg-white hover:shadow-sm border border-transparent rounded-xl transition-all"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                form="driverForm"
                                className="px-8 py-2.5 bg-emerald-600 text-white text-sm font-bold rounded-xl hover:bg-emerald-700 transition-all shadow-lg shadow-emerald-600/20"
                            >
                                {editingDriver ? 'Update Driver' : 'Create Driver'}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Drivers;
