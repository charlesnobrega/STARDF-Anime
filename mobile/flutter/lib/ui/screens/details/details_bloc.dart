import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import '../../../models/anime.dart';
import '../../../core/bridge.dart';

// Events
abstract class DetailsEvent extends Equatable {
  @override
  List<Object> get props => [];
}

class LoadEpisodes extends DetailsEvent {
  final String animeUrl;
  final String source;
  LoadEpisodes(this.animeUrl, this.source);
  @override
  List<Object> get props => [animeUrl, source];
}

// States
abstract class DetailsState extends Equatable {
  @override
  List<Object> get props => [];
}

class DetailsInitial extends DetailsState {}
class DetailsLoading extends DetailsState {}
class DetailsSuccess extends DetailsState {
  final List<Episode> episodes;
  DetailsSuccess(this.episodes);
  @override
  List<Object> get props => [episodes];
}
class DetailsError extends DetailsState {
  final String message;
  DetailsError(this.message);
  @override
  List<Object> get props => [message];
}

// BLoC
class DetailsBloc extends Bloc<DetailsEvent, DetailsState> {
  DetailsBloc() : super(DetailsInitial()) {
    on<LoadEpisodes>((event, emit) async {
      emit(DetailsLoading());
      try {
        final List<dynamic> data = await GoBridge.getEpisodes(event.animeUrl, event.source);
        final List<Episode> eps = data.map((e) => Episode.fromJson(e)).toList();
        emit(DetailsSuccess(eps));
      } catch (e) {
        emit(DetailsError(e.toString()));
      }
    });
  }
}
