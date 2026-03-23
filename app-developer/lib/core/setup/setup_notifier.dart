import 'package:flutter/foundation.dart';

// Notifier untuk status instalasi sistem
class SetupNotifier extends ChangeNotifier {
  bool _isInstalled = true; // default true agar tidak blocking saat API belum ready
  bool _isChecking = true;

  bool get isInstalled => _isInstalled;
  bool get isChecking => _isChecking;

  void setInstalled(bool value) {
    _isInstalled = value;
    _isChecking = false;
    notifyListeners();
  }

  void setChecking(bool value) {
    _isChecking = value;
    notifyListeners();
  }
}
