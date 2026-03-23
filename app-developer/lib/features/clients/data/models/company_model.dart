import '../../domain/entities/company_entity.dart';

class CompanyModel extends CompanyEntity {
  const CompanyModel({
    required super.id,
    required super.code,
    required super.name,
    super.companyType,
    required super.currency,
    required super.isActive,
    super.npwp,
    super.email,
    super.phone,
    super.address,
    super.website,
    required super.modules,
    required super.apps,
    required super.createdAt,
  });

  /// Parsing dari format ClientLicense (api-developer /api/v1/licenses)
  factory CompanyModel.fromJson(Map<String, dynamic> json) {
    return CompanyModel(
      id: json['id'] as String? ?? '',
      code: json['license_key'] as String? ?? '',
      name: json['client_name'] as String? ?? '',
      companyType: json['plan'] as String?,
      currency: 'IDR',
      isActive: (json['status'] as String?) == 'active',
      npwp: null,
      email: json['client_email'] as String?,
      phone: null,
      address: null,
      website: json['flasherp_url'] as String?,
      modules: [],
      apps: [],
      createdAt: DateTime.tryParse(json['created_at'] as String? ?? '') ?? DateTime.now(),
    );
  }
}
