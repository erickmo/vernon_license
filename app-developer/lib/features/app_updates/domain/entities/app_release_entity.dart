class AppReleaseEntity {
  final String id;
  final String appId;
  final String version;
  final int versionCode;
  final String downloadUrl;
  final String releaseNotes;
  final bool isMandatory;
  final DateTime createdAt;

  const AppReleaseEntity({
    required this.id,
    required this.appId,
    required this.version,
    required this.versionCode,
    required this.downloadUrl,
    required this.releaseNotes,
    required this.isMandatory,
    required this.createdAt,
  });
}

class ClientInstallEntity {
  final String id;
  final String companyId;
  final String appId;
  final String installedVersion;
  final int installedVersionCode;
  final String targetVersion;
  final int targetVersionCode;
  final bool forceUpdate;
  final String downloadUrl;
  final String releaseNotes;
  final DateTime? lastCheckAt;
  final DateTime updatedAt;

  const ClientInstallEntity({
    required this.id,
    required this.companyId,
    required this.appId,
    required this.installedVersion,
    required this.installedVersionCode,
    required this.targetVersion,
    required this.targetVersionCode,
    required this.forceUpdate,
    required this.downloadUrl,
    required this.releaseNotes,
    this.lastCheckAt,
    required this.updatedAt,
  });

  bool get needsUpdate => targetVersionCode > installedVersionCode && targetVersionCode > 0;
}
