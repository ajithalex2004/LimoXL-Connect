import { useState, useEffect } from 'react';
import { operatorService } from '../../services/api';
import { Building, Plus, Mail, User, MapPin, Phone, Pencil, Trash2 } from 'lucide-react';

interface OutsourceCompany {
    id: string;
    name: string;
    contact_person?: string;
    email?: string;
    contact_number?: string;
    city?: string;
    country?: string;
    is_active: boolean;
    created_at: string;
}

const OutsourceMaster = () => {
    const [isCreateModalOpen, setCreateModalOpen] = useState(false);
    const [loading, setLoading] = useState(false);
    const [companies, setCompanies] = useState<OutsourceCompany[]>([]);
    const [fetchingCompanies, setFetchingCompanies] = useState(true);
    const [editingCompany, setEditingCompany] = useState<OutsourceCompany | null>(null);
    const [deletingCompanyId, setDeletingCompanyId] = useState<string | null>(null);

    const [formData, setFormData] = useState({
        name: '',
        contact_person: '',
        designation: '',
        email: '',
        contact_number: '',
        address: '',
        city: '',
        country: '',
    });

    useEffect(() => {
        fetchCompanies();
    }, []);

    const fetchCompanies = async () => {
        try {
            setFetchingCompanies(true);
            const data = await operatorService.listOutsourceCompanies();
            setCompanies(data || []);
        } catch (error) {
            console.error('Failed to fetch companies:', error);
        } finally {
            setFetchingCompanies(false);
        }
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        try {
            if (editingCompany) {
                await operatorService.updateOutsourceCompany(editingCompany.id, formData);
                alert('Outsource Company Updated Successfully!');
            } else {
                await operatorService.createOutsourceCompany(formData);
                alert('Outsource Company Created Successfully!');
            }
            setCreateModalOpen(false);
            setEditingCompany(null);
            setFormData({
                name: '',
                contact_person: '',
                designation: '',
                email: '',
                contact_number: '',
                address: '',
                city: '',
                country: '',
            });
            fetchCompanies();
        } catch (error: any) {
            alert('Failed to save company: ' + (error.response?.data?.message || 'Unknown error'));
        } finally {
            setLoading(false);
        }
    };

    const handleEdit = (company: OutsourceCompany) => {
        setEditingCompany(company);
        setFormData({
            name: company.name,
            contact_person: company.contact_person || '',
            designation: '',
            email: company.email || '',
            contact_number: company.contact_number || '',
            address: '',
            city: company.city || '',
            country: company.country || '',
        });
        setCreateModalOpen(true);
    };

    const handleDelete = async (id: string) => {
        if (!window.confirm('Are you sure you want to delete this company?')) return;

        setDeletingCompanyId(id);
        try {
            await operatorService.deleteOutsourceCompany(id);
            alert('Company deleted successfully!');
            fetchCompanies();
        } catch (error: any) {
            alert('Failed to delete company: ' + (error.response?.data?.message || 'Unknown error'));
        } finally {
            setDeletingCompanyId(null);
        }
    };

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">Outsource Companies</h1>
                    <p className="text-gray-500">Manage your limo outsourcing partners.</p>
                </div>
                <button
                    onClick={() => setCreateModalOpen(true)}
                    className="flex items-center gap-2 bg-emerald-600 hover:bg-emerald-700 text-white px-4 py-2 rounded-lg transition-colors"
                >
                    <Plus className="h-5 w-5" />
                    New Company
                </button>
            </div>

            {/* List View */}
            {fetchingCompanies ? (
                <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-8 text-center">
                    <p className="text-gray-500">Loading companies...</p>
                </div>
            ) : companies.length === 0 ? (
                <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-8 text-center">
                    <div className="mx-auto w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-4">
                        <Building className="h-8 w-8 text-gray-400" />
                    </div>
                    <h3 className="text-lg font-medium text-gray-900">No companies found</h3>
                    <p className="text-gray-500 mt-1">Get started by onboarding a new outsource partner.</p>
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {companies.map((company) => (
                        <div key={company.id} className="bg-white rounded-xl shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow">
                            <div className="flex items-start justify-between mb-4">
                                <div className="flex items-center gap-3">
                                    <div className="w-12 h-12 bg-emerald-100 rounded-lg flex items-center justify-center">
                                        <Building className="h-6 w-6 text-emerald-600" />
                                    </div>
                                    <div>
                                        <h3 className="font-semibold text-gray-900">{company.name}</h3>
                                        {company.is_active ? (
                                            <span className="text-xs text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded">Active</span>
                                        ) : (
                                            <span className="text-xs text-gray-600 bg-gray-50 px-2 py-0.5 rounded">Inactive</span>
                                        )}
                                    </div>
                                </div>
                                <div className="flex gap-2">
                                    <button
                                        onClick={() => handleEdit(company)}
                                        className="p-2 text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                                        title="Edit"
                                    >
                                        <Pencil className="h-4 w-4" />
                                    </button>
                                    <button
                                        onClick={() => handleDelete(company.id)}
                                        disabled={deletingCompanyId === company.id}
                                        className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors disabled:opacity-50"
                                        title="Delete"
                                    >
                                        <Trash2 className="h-4 w-4" />
                                    </button>
                                </div>
                            </div>

                            <div className="space-y-2 text-sm">
                                {company.contact_person && (
                                    <div className="flex items-center gap-2 text-gray-600">
                                        <User className="h-4 w-4" />
                                        <span>{company.contact_person}</span>
                                    </div>
                                )}
                                {company.email && (
                                    <div className="flex items-center gap-2 text-gray-600">
                                        <Mail className="h-4 w-4" />
                                        <span className="truncate">{company.email}</span>
                                    </div>
                                )}
                                {company.contact_number && (
                                    <div className="flex items-center gap-2 text-gray-600">
                                        <Phone className="h-4 w-4" />
                                        <span>{company.contact_number}</span>
                                    </div>
                                )}
                                {(company.city || company.country) && (
                                    <div className="flex items-center gap-2 text-gray-600">
                                        <MapPin className="h-4 w-4" />
                                        <span>{[company.city, company.country].filter(Boolean).join(', ')}</span>
                                    </div>
                                )}
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {/* Create Modal */}
            {isCreateModalOpen && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
                    <div className="bg-white rounded-xl shadow-xl max-w-2xl w-full p-6 animate-in fade-in zoom-in duration-200 max-h-[90vh] overflow-y-auto">
                        <div className="flex justify-between items-center mb-6">
                            <h2 className="text-xl font-bold text-gray-900">{editingCompany ? 'Edit Partner' : 'Onboard Partner'}</h2>
                            <button onClick={() => { setCreateModalOpen(false); setEditingCompany(null); }} className="text-gray-400 hover:text-gray-500">
                                <Plus className="h-6 w-6 rotate-45" />
                            </button>
                        </div>

                        <form onSubmit={handleSubmit} className="space-y-4">
                            <div className="grid grid-cols-2 gap-4">
                                <div className="col-span-2">
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Company Name *</label>
                                    <div className="relative">
                                        <Building className="absolute left-3 top-2.5 h-5 w-5 text-gray-400" />
                                        <input
                                            name="name"
                                            type="text"
                                            required
                                            className="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
                                            placeholder="AAT Limo Services"
                                            value={formData.name}
                                            onChange={handleChange}
                                        />
                                    </div>
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Contact Person</label>
                                    <input
                                        name="contact_person"
                                        type="text"
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
                                        placeholder="John Doe"
                                        value={formData.contact_person}
                                        onChange={handleChange}
                                    />
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Designation</label>
                                    <input
                                        name="designation"
                                        type="text"
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
                                        placeholder="Operations Manager"
                                        value={formData.designation}
                                        onChange={handleChange}
                                    />
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
                                    <input
                                        name="email"
                                        type="email"
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
                                        placeholder="contact@company.com"
                                        value={formData.email}
                                        onChange={handleChange}
                                    />
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Contact Number</label>
                                    <input
                                        name="contact_number"
                                        type="tel"
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
                                        placeholder="+971 50 123 4567"
                                        value={formData.contact_number}
                                        onChange={handleChange}
                                    />
                                </div>

                                <div className="col-span-2">
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Address</label>
                                    <input
                                        name="address"
                                        type="text"
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
                                        placeholder="123 Business Bay"
                                        value={formData.address}
                                        onChange={handleChange}
                                    />
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">City</label>
                                    <input
                                        name="city"
                                        type="text"
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
                                        placeholder="Dubai"
                                        value={formData.city}
                                        onChange={handleChange}
                                    />
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Country</label>
                                    <input
                                        name="country"
                                        type="text"
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
                                        placeholder="UAE"
                                        value={formData.country}
                                        onChange={handleChange}
                                    />
                                </div>
                            </div>

                            <button
                                type="submit"
                                disabled={loading}
                                className="w-full mt-4 bg-emerald-600 hover:bg-emerald-700 text-white font-medium py-2.5 px-4 rounded-lg transition-colors shadow-sm disabled:opacity-50"
                            >
                                {loading ? 'Saving...' : (editingCompany ? 'Update Company' : 'Create Company')}
                            </button>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default OutsourceMaster;
