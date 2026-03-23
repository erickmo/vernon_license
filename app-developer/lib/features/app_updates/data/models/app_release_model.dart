import '../../domain/entities/app_release_entity.dart';

class AppReleaseModel extends AppReleaseEntity {
  const AppReleaseModel({
    required super.id,
    required super.appId,
    required super.version,
    required super.versionCode,
    required super.downloadUrl,
    required super.releaseNotes,
    required super.isMandatory,
    required super.createdAt,
  });

  factory AppReleaseModel.fromJson(Map<String, dynamic> json) {
    return AppReleaseModel(
      id: json['id'] as String,
      appId: json['app_id'] as String,
      version: json['version'] as String,
      versionCode: json['version_code'] as int,
      downloadUrl: json['download_url'] as String,
      releaseNotes: json['release_notes'] as String? ?? '',
      isMandatory: json['is_mandatory'] as bool? ?? false,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }
}

class ClientInstallModel extends ClientInstallEntity {
  const ClientInstallModel({
    required super.id,
    required super.companyId,
    required super.appId,
    required super.installedVersion,
    required super.installedVersionCode,
    required super.targetVersion,
    required super.targetVersionCode,
    required super.forceUpdate,
    required super.downloadUrl,
    required super.releaseNotes,
    super.lastCheckAt,
    required super.updatedAt,
  });

  factory ClientInstallModel.fromJson(Map<String, dynamic> json) {
    return ClientInstallModel(
      id: json['id'] as String? ?? '',
      companyId: json['company_id'] as String? ?? '',
      appId: json['app_id'] as String,
      installedVersion: json['installed_version'] as String? ?? '',
      installedVersionCode: json['installed_version_code'] as int? ?? 0,
      targetVersion: json['target_version'] as String? ?? '',
      targetVersionCode: json['target_version_code'] as int? ?? 0,
      forceUpdate: json['force_update'] as bool? ?? false,
      downloadUrl: json['download_url'] as String? ?? '',
      releaseNotes: json['release_notes'] as String? ?? '',
      lastCheckAt: json['last_check_at'] != null
          ? DateTime.tryParse(json['last_check_at'] as String)
          : null,
      updatedAt: DateTime.parse(
          json['updated_at'] as String? ?? DateTime.now().toIso8601String()),
    );
  }
}
