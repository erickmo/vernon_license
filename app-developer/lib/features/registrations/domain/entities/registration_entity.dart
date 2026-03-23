class RegistrationEntity {
  final String id;
  final String nama;
  final String? npwp;
  final String? alamat;
  final String? email;
  final String? telepon;
  final String? catatan;
  final String requestedBy;
  final String status;
  final DateTime createdAt;

  const RegistrationEntity({
    required this.id,
    required this.nama,
    this.npwp,
    this.alamat,
    this.email,
    this.telepon,
    this.catatan,
    required this.requestedBy,
    required this.status,
    required this.createdAt,
  });

  bool get isPending => status == 'pending';
  bool get isApproved => status == 'approved';
  bool get isRejected => status == 'rejected';
}
