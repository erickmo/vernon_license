part of 'setup_cubit.dart';

abstract class SetupState {}

class SetupInitial extends SetupState {}

class SetupLoading extends SetupState {}

class SetupStatusLoaded extends SetupState {
  final bool isInstalled;
  SetupStatusLoaded(this.isInstalled);
}

class SetupInstalling extends SetupState {}

class SetupInstallSuccess extends SetupState {}

class SetupError extends SetupState {
  final String message;
  SetupError(this.message);
}
