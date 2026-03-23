import 'package:flutter/foundation.dart';

/// Menyimpan jumlah pengajuan yang belum diproses.
/// Digunakan untuk badge pada tab Notifikasi di bottom navigation bar.
class PendingCountNotifier extends ValueNotifier<int> {
  PendingCountNotifier() : super(0);

  void update(int count) {
    value = count;
  }
}
