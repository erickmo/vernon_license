import '../../domain/entities/registration_entity.dart';

class RegistrationModel extends RegistrationEntity {
  const RegistrationModel({
    required super.id,
    required super.nama,
    super.npwp,
    super.alamat,
    super.email,
    super.telepon,
    super.catatan,
    required super.requestedBy,
    required super.status,
    required super.createdAt,
  });

  /// Parsing dari format ClientLicense (api-developer /api/v1/licenses)
  factory RegistrationModel.fromJson(Map<String, dynamic> json) {
    return RegistrationModel(
      id:          json['id'] as String? ?? '',
      nama:        json['client_name'] as String? ?? '',
      npwp:        null,
      alamat:      json['flasherp_url'] as String?,
      email:       json['client_email'] as String?,
      telepon:     null,
      catatan:     json['plan'] as String?,
      requestedBy: json['created_by'] as String? ?? '',
      status:      json['status'] as String? ?? 'active',
      createdAt:   DateTime.tryParse(json['created_at'] as String? ?? '') ??
          DateTime.now(),
    );
  }
}
