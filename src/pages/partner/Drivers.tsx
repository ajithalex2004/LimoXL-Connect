import { useState, useEffect } from 'react';
import { partnerService, type Driver } from '../../services/api';
import { Plus, User, Phone, FileText, ShieldCheck } from 'lucide-react';

const Drivers = () => {
    const [drivers, setDrivers] = useState<Driver[]>([]);
    const [loading, setLoading] = useState(true);
    const [showModal, setShowModal] = useState(false);

    // Form State
    const [formData, setFormData] = useState({
        name: '',
        phone: '',
        license_number: '',
        license_expiry: '',
        itc_permit_expiry: '',
        visa_expiry: ''
    });

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        try {
            const data = await partnerService.listDrivers();
            setDrivers(data);
        } catch (error) {
            console.error("Failed to load drivers", error);
        } finally {
            setLoading(false);
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await partnerService.createDriver(formData);
            setShowModal(false);
            setFormData({
                name: '',
                phone: '',
                license_number: '',
                license_expiry: '',
                itc_permit_expiry: '',
                visa_expiry: ''
            });
            loadData();
        } catch (error) {
            console.error("Failed to create driver", error);
            alert("Failed to create driver");
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
                <h1 className="text-2xl font-bold text-gray-900">Drivers</h1>
                <button
                    onClick={() => setShowModal(true)}
                    className="flex items-center gap-2 px-4 py-2 bg-black text-white rounded-lg hover:bg-gray-800"
                >
                    <Plus className="h-4 w-4" /> Add Driver
                </button>
            </div>

            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
                {drivers.map((driver) => (
                    <div key={driver.id} className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
                        <div className="flex items-center gap-4 mb-4">
                            <div className="h-12 w-12 bg-gray-100 rounded-full flex items-center justify-center">
                                <User className="h-6 w-6 text-gray-600" />
                            </div>
                            <div>
                                <h3 className="font-bold text-gray-900">{driver.name}</h3>
                                <p className="text-sm text-green-600 font-medium">Active</p>
                            </div>
                        </div>
                        <div className="space-y-3 text-sm text-gray-600">
                            <div className="flex items-center gap-2">
                                <Phone className="h-4 w-4 text-gray-400" /> {driver.phone}
                            </div>
                            <div className="flex items-center gap-2">
                                <FileText className="h-4 w-4 text-gray-400" /> {driver.license_number}
                            </div>

                            <div className="pt-2 border-t border-gray-100 space-y-2">
                                <p className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Compliance Status</p>
                                <div className="grid grid-cols-1 gap-1.5">
                                    <div className="flex items-center justify-between">
                                        <span className="text-xs">Driving License</span>
                                        {(() => {
                                            const status = getExpiryStatus(driver.license_expiry);
                                            return status ? (
                                                <span className={`px-2 py-0.5 rounded text-[10px] font-bold border ${status.color}`}>
                                                    {status.label}
                                                </span>
                                            ) : <span className="text-gray-400 text-[10px]">Not Provided</span>
                                        })()}
                                    </div>
                                    <div className="flex items-center justify-between">
                                        <span className="text-xs">ITC Permit</span>
                                        {(() => {
                                            const status = getExpiryStatus(driver.itc_permit_expiry);
                                            return status ? (
                                                <span className={`px-2 py-0.5 rounded text-[10px] font-bold border ${status.color}`}>
                                                    {status.label}
                                                </span>
                                            ) : <span className="text-gray-400 text-[10px]">Not Provided</span>
                                        })()}
                                    </div>
                                    <div className="flex items-center justify-between">
                                        <span className="text-xs">Visa Status</span>
                                        {(() => {
                                            const status = getExpiryStatus(driver.visa_expiry);
                                            return status ? (
                                                <span className={`px-2 py-0.5 rounded text-[10px] font-bold border ${status.color}`}>
                                                    {status.label}
                                                </span>
                                            ) : <span className="text-gray-400 text-[10px]">Not Provided</span>
                                        })()}
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                ))}
            </div>

            {showModal && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
                    <div className="bg-white rounded-xl shadow-xl max-w-md w-full p-6">
                        <h2 className="text-xl font-bold mb-4">Add New Driver</h2>
                        <form onSubmit={handleSubmit} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
                                <input
                                    type="text"
                                    required
                                    className="w-full px-3 py-2 border rounded-lg"
                                    value={formData.name}
                                    onChange={e => setFormData({ ...formData, name: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Phone Number</label>
                                <input
                                    type="tel"
                                    required
                                    className="w-full px-3 py-2 border rounded-lg"
                                    value={formData.phone}
                                    onChange={e => setFormData({ ...formData, phone: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">License Number</label>
                                <input
                                    type="text"
                                    required
                                    className="w-full px-3 py-2 border rounded-lg"
                                    value={formData.license_number}
                                    onChange={e => setFormData({ ...formData, license_number: e.target.value })}
                                />
                            </div>

                            <div className="bg-gray-50 p-4 rounded-lg space-y-4">
                                <p className="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-2">
                                    <ShieldCheck className="h-3 w-3" /> Compliance Verification
                                </p>
                                <div className="grid grid-cols-1 gap-4">
                                    <div>
                                        <label className="block text-[11px] font-bold text-gray-600 mb-1 uppercase">License Expiry</label>
                                        <input
                                            type="date"
                                            className="w-full px-3 py-1.5 border rounded-lg text-sm"
                                            value={formData.license_expiry}
                                            onChange={e => setFormData({ ...formData, license_expiry: e.target.value })}
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-[11px] font-bold text-gray-600 mb-1 uppercase">ITC Permit Expiry</label>
                                        <input
                                            type="date"
                                            className="w-full px-3 py-1.5 border rounded-lg text-sm"
                                            value={formData.itc_permit_expiry}
                                            onChange={e => setFormData({ ...formData, itc_permit_expiry: e.target.value })}
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-[11px] font-bold text-gray-600 mb-1 uppercase">Visa Expiry</label>
                                        <input
                                            type="date"
                                            className="w-full px-3 py-1.5 border rounded-lg text-sm"
                                            value={formData.visa_expiry}
                                            onChange={e => setFormData({ ...formData, visa_expiry: e.target.value })}
                                        />
                                    </div>
                                </div>
                            </div>
                            <div className="flex gap-3 pt-4">
                                <button type="button" onClick={() => setShowModal(false)} className="flex-1 py-2 border rounded-lg hover:bg-gray-50">Cancel</button>
                                <button type="submit" className="flex-1 py-2 bg-black text-white rounded-lg hover:bg-gray-800">Add Driver</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Drivers;
