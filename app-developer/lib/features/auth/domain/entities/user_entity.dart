class UserEntity {
  final String id;
  final String name;
  final String email;
  final String role;

  const UserEntity({
    required this.id,
    required this.name,
    required this.email,
    required this.role,
  });

  bool get isSuperuser => role == 'superuser';
  bool get isFinance => role == 'finance';
  bool get isProjectOwner => role == 'project_owner';
  bool get isProjectManager => role == 'project_manager';

  // superuser dapat akses semua fitur
  bool get canManageRegistrations => isSuperuser || isProjectOwner;
  bool get canManageLicense => isSuperuser || isFinance;
  bool get canManageClientSettings => isSuperuser || isProjectManager;
  bool get canUpdatePolicy => isSuperuser || isProjectOwner;
}
