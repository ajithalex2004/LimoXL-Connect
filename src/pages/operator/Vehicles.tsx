import { useState, useEffect } from 'react';
import { operatorService, type Vehicle, type FleetAttachment, type NUIMaster } from '../../services/api';
import {
    Plus, Car, Search, Edit2, Trash2,
    X, Info, Paperclip, CheckCircle2,
    Calendar, Shield, MapPin, Hash, Palette,
    Download, FileText
} from 'lucide-react';

const Vehicles = () => {
    // List State
    const [vehicles, setVehicles] = useState<Vehicle[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');

    // Masters State
    const [vehicleTypes, setVehicleTypes] = useState<NUIMaster[]>([]);
    const [vehicleClasses, setVehicleClasses] = useState<NUIMaster[]>([]);
    const [vehicleUsages, setVehicleUsages] = useState<NUIMaster[]>([]);
    const [hierarchies, setHierarchies] = useState<NUIMaster[]>([]);

    // Modal State
    const [showModal, setShowModal] = useState(false);
    const [activeTab, setActiveTab] = useState<'basic' | 'attachments'>('basic');
    const [editingVehicle, setEditingVehicle] = useState<Vehicle | null>(null);
    const [attachments, setAttachments] = useState<FleetAttachment[]>([]);
    const [attLoading, setAttLoading] = useState(false);

    // Form State
    const [formData, setFormData] = useState<Partial<Vehicle>>({
        plate_number: '',
        model: '',
        color: '',
        vehicle_class: '',
        vehicle_group: '',
        capacity: 4,
        year_of_manufacture: new Date().getFullYear(),
        vin: '',
        chassis_no: '',
        registration_number: '',
        plate_code: '',
        plate_category: '',
        emirate: 'Dubai',
        hierarchy: '',
        vehicle_usage: '',
        permit_expiry: '',
        insurance_expiry: '',
        status: 'OFFLINE'
    });

    useEffect(() => {
        loadData();
        loadMasters();
    }, []);

    const loadMasters = async () => {
        try {
            console.log("DEBUG: Loading NUI Masters for Vehicles...");
            const [types, classes, usages, hierarchyList] = await Promise.all([
                operatorService.listMasters('TYPE'),
                operatorService.listMasters('CLASS'),
                operatorService.listMasters('USAGE'),
                operatorService.listMasters('HIERARCHY')
            ]);
            console.log("DEBUG: Fetched Masters:", {
                types: types?.length || 0,
                classes: classes?.length || 0,
                usages: usages?.length || 0,
                hierarchies: hierarchyList?.length || 0
            });
            setVehicleTypes(types?.filter(m => m.is_active) || []);
            setVehicleClasses(classes?.filter(m => m.is_active) || []);
            setVehicleUsages(usages?.filter(m => m.is_active) || []);
            setHierarchies(hierarchyList?.filter(m => m.is_active) || []);
        } catch (error) {
            console.error("Failed to load masters", error);
        }
    };

    const loadData = async () => {
        try {
            setLoading(true);
            const data = await operatorService.listVehicles();
            setVehicles(data || []);
        } catch (error) {
            console.error("Failed to load vehicles", error);
        } finally {
            setLoading(false);
        }
    };

    const handleOpenAdd = () => {
        loadMasters(); // Reload to get latest data
        setEditingVehicle(null);
        setFormData({
            plate_number: '',
            model: '',
            color: '',
            vehicle_class: 'Standard',
            vehicle_group: 'Sedan',
            capacity: 4,
            year_of_manufacture: new Date().getFullYear(),
            vin: '',
            chassis_no: '',
            registration_number: '',
            plate_code: '',
            plate_category: '',
            emirate: 'Dubai',
            hierarchy: '',
            vehicle_usage: 'Public',
            permit_expiry: '',
            insurance_expiry: '',
            status: 'OFFLINE'
        });
        setActiveTab('basic');
        setShowModal(true);
    };

    const handleOpenEdit = async (vehicle: Vehicle) => {
        setEditingVehicle(vehicle);
        setFormData({ ...vehicle });
        setActiveTab('basic');
        setShowModal(true);

        // Load attachments
        try {
            setAttLoading(true);
            const atts = await operatorService.listAttachments(vehicle.id, 'vehicle');
            setAttachments(atts || []);
        } catch (error) {
            console.error("Failed to load attachments", error);
        } finally {
            setAttLoading(false);
        }
    };

    const handleDelete = async (id: string) => {
        if (!confirm('Are you sure you want to delete this vehicle?')) return;
        try {
            await operatorService.deleteVehicle(id);
            loadData();
        } catch (error) {
            console.error("Failed to delete vehicle", error);
            alert("Failed to delete vehicle");
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            if (editingVehicle) {
                await operatorService.updateVehicle(editingVehicle.id, formData);
            } else {
                await operatorService.createVehicle(formData);
            }
            setShowModal(false);
            loadData();
        } catch (error) {
            console.error("Failed to save vehicle", error);
            alert("Failed to save vehicle");
        }
    };

    const filteredVehicles = vehicles.filter(v =>
        v.plate_number.toLowerCase().includes(searchTerm.toLowerCase()) ||
        v.model.toLowerCase().includes(searchTerm.toLowerCase())
    );

    if (loading && vehicles.length === 0) {
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
                    <h1 className="text-2xl font-bold text-gray-900">Vehicle Master</h1>
                    <p className="text-gray-500 text-sm">Manage your fleet vehicles and registration details</p>
                </div>
                <div className="flex items-center gap-3">
                    <div className="relative">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
                        <input
                            type="text"
                            placeholder="Search by plate or model..."
                            className="pl-10 pr-4 py-2 border border-gray-200 rounded-lg text-sm w-64 focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500"
                            value={searchTerm}
                            onChange={e => setSearchTerm(e.target.value)}
                        />
                    </div>
                    <button
                        onClick={handleOpenAdd}
                        className="flex items-center gap-2 px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors text-sm font-semibold shadow-sm"
                    >
                        <Plus className="h-4 w-4" /> Add Vehicle
                    </button>
                </div>
            </div>

            {/* Vehicle List */}
            <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full text-left">
                        <thead className="bg-gray-50 border-b border-gray-200">
                            <tr>
                                <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase">Vehicle Info</th>
                                <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase">Registration</th>
                                <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase">Class/Group</th>
                                <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase">Status</th>
                                <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase">Compliance</th>
                                <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100">
                            {filteredVehicles.length === 0 ? (
                                <tr>
                                    <td colSpan={6} className="px-6 py-12 text-center text-gray-500">
                                        No vehicles found. Add your first vehicle to get started.
                                    </td>
                                </tr>
                            ) : (
                                filteredVehicles.map((vehicle) => (
                                    <tr key={vehicle.id} className="hover:bg-gray-50/50 transition-colors group">
                                        <td className="px-6 py-4">
                                            <div className="flex items-center gap-3">
                                                <div className="h-10 w-10 bg-emerald-50 rounded-lg flex items-center justify-center text-emerald-600">
                                                    <Car className="h-5 w-5" />
                                                </div>
                                                <div>
                                                    <div className="font-bold text-gray-900">{vehicle.plate_number}</div>
                                                    <div className="text-xs text-gray-500 uppercase">{vehicle.model} • {vehicle.color}</div>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className="text-sm font-medium text-gray-700">{vehicle.registration_number || '-'}</div>
                                            <div className="text-xs text-gray-400">{vehicle.emirate} • {vehicle.plate_code}</div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className="text-sm text-gray-700">{vehicle.vehicle_class}</div>
                                            <div className="text-xs text-gray-400">{vehicle.vehicle_group}</div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <span className={`px-2.5 py-1 text-[10px] font-bold rounded-full uppercase border ${vehicle.status === 'IDLE' ? 'bg-emerald-50 text-emerald-700 border-emerald-100' :
                                                vehicle.status === 'ON_TRIP' ? 'bg-blue-50 text-blue-700 border-blue-100' :
                                                    'bg-gray-50 text-gray-600 border-gray-100'
                                                }`}>
                                                {vehicle.status}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className="space-y-1">
                                                <div className="flex items-center gap-2">
                                                    <div className={`h-1.5 w-1.5 rounded-full ${new Date(vehicle.permit_expiry || '') > new Date() ? 'bg-emerald-500' : 'bg-red-500'}`}></div>
                                                    <span className="text-[10px] text-gray-500 uppercase">Permit: {vehicle.permit_expiry ? new Date(vehicle.permit_expiry).toLocaleDateString() : 'N/A'}</span>
                                                </div>
                                                <div className="flex items-center gap-2">
                                                    <div className={`h-1.5 w-1.5 rounded-full ${new Date(vehicle.insurance_expiry || '') > new Date() ? 'bg-emerald-500' : 'bg-red-500'}`}></div>
                                                    <span className="text-[10px] text-gray-500 uppercase">Insurance: {vehicle.insurance_expiry ? new Date(vehicle.insurance_expiry).toLocaleDateString() : 'N/A'}</span>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4 text-right">
                                            <div className="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                                <button onClick={() => handleOpenEdit(vehicle)} className="p-1.5 text-gray-400 hover:text-emerald-600 hover:bg-emerald-50 rounded-lg">
                                                    <Edit2 className="h-4 w-4" />
                                                </button>
                                                <button onClick={() => handleDelete(vehicle.id)} className="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg">
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
                            <h2 className="text-xl font-bold text-gray-900">{editingVehicle ? 'Edit Vehicle' : 'Create New Vehicle'}</h2>
                            <button onClick={() => setShowModal(false)} className="p-2 hover:bg-gray-100 rounded-lg transition-colors">
                                <X className="h-5 w-5 text-gray-500" />
                            </button>
                        </div>

                        {/* Tabs */}
                        <div className="px-6 flex border-b border-gray-100">
                            <button
                                onClick={() => setActiveTab('basic')}
                                className={`px-4 py-3 text-sm font-semibold flex items-center gap-2 border-b-2 transition-colors ${activeTab === 'basic' ? 'border-emerald-600 text-emerald-600' : 'border-transparent text-gray-500 hover:text-gray-700'
                                    }`}
                            >
                                <Info className="h-4 w-4" /> Basic Details
                            </button>
                            <button
                                onClick={() => setActiveTab('attachments')}
                                className={`px-4 py-3 text-sm font-semibold flex items-center gap-2 border-b-2 transition-colors ${activeTab === 'attachments' ? 'border-emerald-600 text-emerald-600' : 'border-transparent text-gray-500 hover:text-gray-700'
                                    }`}
                            >
                                <Paperclip className="h-4 w-4" /> Attachments
                            </button>
                        </div>

                        {/* Modal Body */}
                        <div className="flex-1 overflow-y-auto p-6">
                            <form id="vehicleForm" onSubmit={handleSubmit} className="space-y-8">
                                {activeTab === 'basic' ? (
                                    <>
                                        {/* General Info */}
                                        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                                            <div className="space-y-1.5">
                                                <label className="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                                                    <Hash className="h-3 w-3" /> Vehicle ID (Plate)
                                                </label>
                                                <input
                                                    type="text"
                                                    required
                                                    placeholder="e.g., VEH001"
                                                    className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                    value={formData.plate_number}
                                                    onChange={e => setFormData({ ...formData, plate_number: e.target.value })}
                                                />
                                            </div>
                                            <div className="space-y-1.5">
                                                <label className="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                                                    <Shield className="h-3 w-3" /> VIN
                                                </label>
                                                <input
                                                    type="text"
                                                    placeholder="Enter VIN"
                                                    className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                    value={formData.vin}
                                                    onChange={e => setFormData({ ...formData, vin: e.target.value })}
                                                />
                                            </div>
                                            <div className="space-y-1.5">
                                                <label className="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                                                    <Info className="h-3 w-3" /> Chassis No
                                                </label>
                                                <input
                                                    type="text"
                                                    placeholder="Enter chassis number"
                                                    className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                    value={formData.chassis_no}
                                                    onChange={e => setFormData({ ...formData, chassis_no: e.target.value })}
                                                />
                                            </div>
                                            <div className="space-y-1.5">
                                                <label className="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                                                    <Calendar className="h-3 w-3" /> Year of Manufacture
                                                </label>
                                                <input
                                                    type="number"
                                                    placeholder="e.g., 2023"
                                                    className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                    value={formData.year_of_manufacture}
                                                    onChange={e => setFormData({ ...formData, year_of_manufacture: parseInt(e.target.value) })}
                                                />
                                            </div>
                                            <div className="space-y-1.5 md:col-span-2">
                                                <label className="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                                                    <Palette className="h-3 w-3" /> Color
                                                </label>
                                                <select
                                                    className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm appearance-none"
                                                    value={formData.color}
                                                    onChange={e => setFormData({ ...formData, color: e.target.value })}
                                                >
                                                    <option value="">Select vehicle color</option>
                                                    <option value="White">White</option>
                                                    <option value="Black">Black</option>
                                                    <option value="Silver">Silver</option>
                                                    <option value="Grey">Grey</option>
                                                    <option value="Gold">Gold</option>
                                                </select>
                                            </div>
                                        </div>

                                        {/* Registration Details Section */}
                                        <div className="pt-4 space-y-4">
                                            <div className="flex items-center gap-2 pb-2 border-b border-gray-100">
                                                <Shield className="h-4 w-4 text-emerald-600" />
                                                <span className="text-sm font-bold text-gray-900">Registration Details</span>
                                            </div>
                                            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                                                <div className="space-y-1.5">
                                                    <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Registration Number</label>
                                                    <input
                                                        type="text"
                                                        placeholder="Enter reg number"
                                                        className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                        value={formData.registration_number}
                                                        onChange={e => setFormData({ ...formData, registration_number: e.target.value })}
                                                    />
                                                </div>
                                                <div className="space-y-1.5">
                                                    <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Plate Code</label>
                                                    <input
                                                        type="text"
                                                        placeholder="e.g., A"
                                                        className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                        value={formData.plate_code}
                                                        onChange={e => setFormData({ ...formData, plate_code: e.target.value })}
                                                    />
                                                </div>
                                                <div className="space-y-1.5">
                                                    <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Emirate</label>
                                                    <select
                                                        className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm appearance-none"
                                                        value={formData.emirate}
                                                        onChange={e => setFormData({ ...formData, emirate: e.target.value })}
                                                    >
                                                        <option value="Dubai">Dubai</option>
                                                        <option value="Abu Dhabi">Abu Dhabi</option>
                                                        <option value="Sharjah">Sharjah</option>
                                                        <option value="Ajman">Ajman</option>
                                                    </select>
                                                </div>
                                            </div>
                                        </div>

                                        {/* Fleet Assignment Section */}
                                        <div className="pt-4 space-y-4">
                                            <div className="flex items-center gap-2 pb-2 border-b border-gray-100">
                                                <MapPin className="h-4 w-4 text-emerald-600" />
                                                <span className="text-sm font-bold text-gray-900">Fleet Assignment</span>
                                            </div>
                                            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                                                <div className="space-y-1.5">
                                                    <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Hierarchy</label>
                                                    <select
                                                        className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm appearance-none"
                                                        value={formData.hierarchy}
                                                        onChange={e => setFormData({ ...formData, hierarchy: e.target.value })}
                                                    >
                                                        <option value="">Select hierarchy</option>
                                                        {hierarchies.map(m => (
                                                            <option key={m.id} value={m.name}>{m.name}</option>
                                                        ))}
                                                    </select>
                                                </div>
                                                <div className="space-y-1.5">
                                                    <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Vehicle Type</label>
                                                    <select
                                                        className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm appearance-none"
                                                        value={formData.vehicle_group}
                                                        onChange={e => setFormData({ ...formData, vehicle_group: e.target.value })}
                                                    >
                                                        <option value="">Select Type</option>
                                                        {vehicleTypes.map(m => (
                                                            <option key={m.id} value={m.name}>{m.name}</option>
                                                        ))}
                                                    </select>
                                                </div>
                                                <div className="space-y-1.5">
                                                    <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Vehicle Class</label>
                                                    <select
                                                        className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm appearance-none"
                                                        value={formData.vehicle_class}
                                                        onChange={e => setFormData({ ...formData, vehicle_class: e.target.value })}
                                                    >
                                                        <option value="">Select Class</option>
                                                        {vehicleClasses.map(m => (
                                                            <option key={m.id} value={m.name}>{m.name}</option>
                                                        ))}
                                                    </select>
                                                </div>
                                                <div className="space-y-1.5">
                                                    <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Vehicle Usage</label>
                                                    <select
                                                        className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm appearance-none"
                                                        value={formData.vehicle_usage}
                                                        onChange={e => setFormData({ ...formData, vehicle_usage: e.target.value })}
                                                    >
                                                        <option value="">Select Usage</option>
                                                        {vehicleUsages.map(m => (
                                                            <option key={m.id} value={m.name}>{m.name}</option>
                                                        ))}
                                                    </select>
                                                </div>
                                            </div>
                                        </div>

                                        {/* Compliance Section */}
                                        <div className="pt-4 space-y-4">
                                            <div className="flex items-center gap-2 pb-2 border-b border-gray-100">
                                                <CheckCircle2 className="h-4 w-4 text-emerald-600" />
                                                <span className="text-sm font-bold text-gray-900">Compliance & Expiry</span>
                                            </div>
                                            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                                <div className="space-y-1.5">
                                                    <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Permit Expiry</label>
                                                    <input
                                                        type="date"
                                                        className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                        value={formData.permit_expiry ? formData.permit_expiry.split('T')[0] : ''}
                                                        onChange={e => setFormData({ ...formData, permit_expiry: e.target.value })}
                                                    />
                                                </div>
                                                <div className="space-y-1.5">
                                                    <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Insurance Expiry</label>
                                                    <input
                                                        type="date"
                                                        className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                                        value={formData.insurance_expiry ? formData.insurance_expiry.split('T')[0] : ''}
                                                        onChange={e => setFormData({ ...formData, insurance_expiry: e.target.value })}
                                                    />
                                                </div>
                                            </div>
                                        </div>
                                    </>
                                ) : (
                                    <div className="space-y-6 animate-in slide-in-from-right-2 duration-300">
                                        {!editingVehicle ? (
                                            <div className="p-12 text-center text-gray-500 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200">
                                                <Info className="mx-auto h-8 w-8 mb-2 opacity-20" />
                                                <p>Please save the vehicle first before adding attachments.</p>
                                            </div>
                                        ) : (
                                            <>
                                                <div className="p-12 border-2 border-dashed border-gray-200 rounded-2xl flex flex-col items-center justify-center gap-4 bg-gray-50/50">
                                                    <div className="p-4 bg-white rounded-2xl shadow-sm border border-gray-100">
                                                        <Paperclip className="h-8 w-8 text-gray-400" />
                                                    </div>
                                                    <div className="text-center">
                                                        <p className="font-bold text-gray-900">Upload Attachments</p>
                                                        <p className="text-xs text-gray-500 mt-1">Mulkiya, Insurance Policy, Emirates ID etc. (PDF, JPG, PNG)</p>
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
                                                                                if (confirm('Delete attachment?')) {
                                                                                    await operatorService.deleteAttachment(att.id);
                                                                                    const updated = await operatorService.listAttachments(editingVehicle.id, 'vehicle');
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
                                form="vehicleForm"
                                className="px-8 py-2.5 bg-emerald-600 text-white text-sm font-bold rounded-xl hover:bg-emerald-700 transition-all shadow-lg shadow-emerald-600/20"
                            >
                                {editingVehicle ? 'Update Vehicle' : 'Create Vehicle'}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Vehicles;
