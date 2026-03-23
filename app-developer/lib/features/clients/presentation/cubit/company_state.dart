part of 'company_cubit.dart';

abstract class CompanyState {}

class CompanyInitial extends CompanyState {}

class CompanyLoading extends CompanyState {}

class CompanyLoaded extends CompanyState {
  final List<CompanyEntity> items;
  final String search;
  CompanyLoaded(this.items, {this.search = ''});
}

class CompanyEmpty extends CompanyState {
  final String search;
  CompanyEmpty({this.search = ''});
}

class CompanyError extends CompanyState {
  final String message;
  CompanyError(this.message);
}
