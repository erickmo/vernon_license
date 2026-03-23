class CompanyEntity {
  final String id;
  final String code;
  final String name;
  final String? companyType;
  final String currency;
  final bool isActive;
  final String? npwp;
  final String? email;
  final String? phone;
  final String? address;
  final String? website;
  final List<String> modules;
  final List<String> apps;
  final DateTime createdAt;

  const CompanyEntity({
    required this.id,
    required this.code,
    required this.name,
    this.companyType,
    required this.currency,
    required this.isActive,
    this.npwp,
    this.email,
    this.phone,
    this.address,
    this.website,
    required this.modules,
    required this.apps,
    required this.createdAt,
  });
}
