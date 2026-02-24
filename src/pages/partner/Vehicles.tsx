import { useState, useEffect } from 'react';
import { partnerService, type Vehicle } from '../../services/api';
import { Plus, Car, ShieldCheck } from 'lucide-react';

const Vehicles = () => {
    // List State
    const [vehicles, setVehicles] = useState<Vehicle[]>([]);
    const [loading, setLoading] = useState(true);
    const [showModal, setShowModal] = useState(false);

    // Form State
    const [formData, setFormData] = useState({
        plate_number: '',
        vehicle_class: 'Standard',
        vehicle_group: 'Sedan',
        model: '',
        capacity: 4 as number,
        permit_expiry: '',
        insurance_expiry: ''
    });

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        try {
            const data = await partnerService.listVehicles();
            setVehicles(data);
        } catch (error) {
            console.error("Failed to load vehicles", error);
        } finally {
            setLoading(false);
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await partnerService.createVehicle(formData);
            setShowModal(false);
            setFormData({
                plate_number: '',
                vehicle_class: 'Standard',
                vehicle_group: 'Sedan',
                model: '',
                capacity: 4,
                permit_expiry: '',
                insurance_expiry: ''
            });
            loadData();
        } catch (error) {
            console.error("Failed to create vehicle", error);
            alert("Failed to create vehicle");
        }
    };

    const getExpiryStatus = (dateStr?: string) => {
        if (!dateStr) return null;
        const expiry = new Date(dateStr);
        const today = new Date();
        const diffDays = Math.ceil((expiry.getTime() - today.getTime()) / (1000 * 3600 * 24));

        if (diffDays < 0) return { label: 'Expired', color: 'text-red-600 bg-red-50 border-red-100' };
        if (diffDays < 30) return { label: `Expires in ${diffDays}d`, color: 'text-amber-600 bg-amber-50 border-amber-100' };
        return { label: 'Compliant', color: 'text-emerald-600 bg-emerald-50 border-emerald-100' };
    };

    if (loading) return <div>Loading...</div>;

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-2xl font-bold text-gray-900">Fleet Management</h1>
                <button
                    onClick={() => setShowModal(true)}
                    className="flex items-center gap-2 px-4 py-2 bg-black text-white rounded-lg hover:bg-gray-800"
                >
                    <Plus className="h-4 w-4" /> Add Vehicle
                </button>
            </div>

            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
                {vehicles.map((vehicle) => (
                    <div key={vehicle.id} className="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
                        <div className="flex justify-between items-start mb-4">
                            <div className="p-3 bg-gray-50 rounded-lg">
                                <Car className="h-6 w-6 text-gray-700" />
                            </div>
                            <span className={`px-2 py-1 text-xs font-medium rounded-full ${vehicle.status === 'IDLE' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'
                                }`}>
                                {vehicle.status}
                            </span>
                        </div>
                        <h3 className="text-lg font-bold text-gray-900 mb-1">{vehicle.plate_number}</h3>
                        <p className="text-sm text-gray-500">
                            {vehicle.vehicle_class} {vehicle.vehicle_group}
                        </p>
                        <p className="text-xs text-gray-400 mb-3">
                            {vehicle.model} • {vehicle.capacity} Seats
                        </p>

                        <div className="pt-3 border-t border-gray-100 space-y-2">
                            <div className="flex items-center justify-between">
                                <span className="text-xs text-gray-500">Permit Expiry</span>
                                {(() => {
                                    const status = getExpiryStatus(vehicle.permit_expiry);
                                    return status ? (
                                        <span className={`px-2 py-0.5 rounded text-[10px] font-bold border ${status.color}`}>
                                            {status.label}
                                        </span>
                                    ) : <span className="text-gray-400 text-[10px]">Pending</span>
                                })()}
                            </div>
                            <div className="flex items-center justify-between">
                                <span className="text-xs text-gray-500">Insurance Expiry</span>
                                {(() => {
                                    const status = getExpiryStatus(vehicle.insurance_expiry);
                                    return status ? (
                                        <span className={`px-2 py-0.5 rounded text-[10px] font-bold border ${status.color}`}>
                                            {status.label}
                                        </span>
                                    ) : <span className="text-gray-400 text-[10px]">Pending</span>
                                })()}
                            </div>
                        </div>
                    </div>
                ))}
            </div>

            {showModal && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
                    <div className="bg-white rounded-xl shadow-xl max-w-md w-full p-6">
                        <h2 className="text-xl font-bold mb-4">Add New Vehicle</h2>
                        <form onSubmit={handleSubmit} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Plate Number</label>
                                <input
                                    type="text"
                                    required
                                    className="w-full px-3 py-2 border rounded-lg"
                                    value={formData.plate_number}
                                    onChange={e => setFormData({ ...formData, plate_number: e.target.value })}
                                />
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Class</label>
                                    <select
                                        className="w-full px-3 py-2 border rounded-lg"
                                        value={formData.vehicle_class}
                                        onChange={e => setFormData({ ...formData, vehicle_class: e.target.value })}
                                    >
                                        <option value="Standard">Standard</option>
                                        <option value="Premium">Premium</option>
                                        <option value="Luxury">Luxury</option>
                                    </select>
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Group</label>
                                    <select
                                        className="w-full px-3 py-2 border rounded-lg"
                                        value={formData.vehicle_group}
                                        onChange={e => setFormData({ ...formData, vehicle_group: e.target.value })}
                                    >
                                        <option value="Sedan">Sedan</option>
                                        <option value="SUV">SUV</option>
                                        <option value="Van">Van</option>
                                        <option value="Bus">Bus</option>
                                    </select>
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Model (Type)</label>
                                <input
                                    type="text"
                                    placeholder="e.g. Lexus ES350"
                                    className="w-full px-3 py-2 border rounded-lg"
                                    value={formData.model}
                                    onChange={e => setFormData({ ...formData, model: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Capacity</label>
                                <input
                                    type="number"
                                    required
                                    min="1"
                                    className="w-full px-3 py-2 border rounded-lg"
                                    value={formData.capacity}
                                    onChange={e => setFormData({ ...formData, capacity: parseInt(e.target.value) })}
                                />
                            </div>

                            <div className="bg-gray-50 p-4 rounded-lg space-y-4">
                                <p className="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-2">
                                    <ShieldCheck className="h-3 w-3" /> Vehicle Compliance
                                </p>
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <label className="block text-[11px] font-bold text-gray-600 mb-1 uppercase">Permit Expiry</label>
                                        <input
                                            type="date"
                                            className="w-full px-3 py-1.5 border rounded-lg text-sm"
                                            value={formData.permit_expiry}
                                            onChange={e => setFormData({ ...formData, permit_expiry: e.target.value })}
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-[11px] font-bold text-gray-600 mb-1 uppercase">Insurance Expiry</label>
                                        <input
                                            type="date"
                                            className="w-full px-3 py-1.5 border rounded-lg text-sm"
                                            value={formData.insurance_expiry}
                                            onChange={e => setFormData({ ...formData, insurance_expiry: e.target.value })}
                                        />
                                    </div>
                                </div>
                            </div>
                            <div className="flex gap-3 pt-4">
                                <button type="button" onClick={() => setShowModal(false)} className="flex-1 py-2 border rounded-lg hover:bg-gray-50">Cancel</button>
                                <button type="submit" className="flex-1 py-2 bg-black text-white rounded-lg hover:bg-gray-800">Add Vehicle</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Vehicles;
