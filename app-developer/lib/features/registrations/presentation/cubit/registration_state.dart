part of 'registration_cubit.dart';

abstract class RegistrationState {}

class RegistrationInitial extends RegistrationState {}

class RegistrationLoading extends RegistrationState {}

class RegistrationLoaded extends RegistrationState {
  final List<RegistrationEntity> items;
  final String? activeFilter;
  RegistrationLoaded(this.items, {this.activeFilter});
}

class RegistrationEmpty extends RegistrationState {
  final String? activeFilter;
  RegistrationEmpty({this.activeFilter});
}

class RegistrationError extends RegistrationState {
  final String message;
  RegistrationError(this.message);
}

class RegistrationActionLoading extends RegistrationState {
  final List<RegistrationEntity> items;
  final String? activeFilter;
  RegistrationActionLoading(this.items, {this.activeFilter});
}

class RegistrationActionSuccess extends RegistrationState {
  final String message;
  final List<RegistrationEntity> items;
  final String? activeFilter;
  RegistrationActionSuccess(this.message, this.items, {this.activeFilter});
}

class ClientCreateLoading extends RegistrationState {
  final List<RegistrationEntity> items;
  final String? activeFilter;
  ClientCreateLoading(this.items, {this.activeFilter});
}

class ClientCreateSuccess extends RegistrationState {
  final String message;
  final List<RegistrationEntity> items;
  final String? activeFilter;
  ClientCreateSuccess(this.message, this.items, {this.activeFilter});
}

class ClientCreateError extends RegistrationState {
  final String message;
  final List<RegistrationEntity> items;
  final String? activeFilter;
  ClientCreateError(this.message, this.items, {this.activeFilter});
}
