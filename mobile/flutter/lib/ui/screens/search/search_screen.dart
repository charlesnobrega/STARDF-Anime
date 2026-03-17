import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:shimmer/shimmer.dart';
import 'package:cached_network_image/cached_network_image.dart';
import '../../../core/app_colors.dart';
import '../../widgets/glass_card.dart';
import 'search_bloc.dart';
import '../../../models/anime.dart';
import '../details/details_screen.dart';

class SearchScreen extends StatelessWidget {
  const SearchScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (context) => SearchBloc(),
      child: Scaffold(
        backgroundColor: AppColors.background,
        body: SafeArea(
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildHeader(),
                const SizedBox(height: 25),
                _buildSearchInput(context),
                const SizedBox(height: 25),
                _buildResultsList(),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildHeader() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          "Explorar",
          style: GoogleFonts.outfit(
            fontSize: 32,
            fontWeight: FontWeight.bold,
            color: AppColors.textMain,
          ),
        ),
        Text(
          "Busque por seus animes ou filmes favoritos",
          style: GoogleFonts.outfit(
            fontSize: 16,
            color: AppColors.textMuted,
          ),
        ),
      ],
    );
  }

  Widget _buildSearchInput(BuildContext context) {
    return GlassCard(
      height: 60,
      padding: const EdgeInsets.symmetric(horizontal: 15),
      child: TextField(
        onChanged: (value) {
          // Debouncing would be better, but for now:
          context.read<SearchBloc>().add(SearchQueryChanged(value));
        },
        style: const TextStyle(color: AppColors.textMain),
        decoration: InputDecoration(
          icon: const Icon(Icons.search, color: AppColors.accent),
          hintText: "Título, gênero ou ano...",
          hintStyle: TextStyle(color: AppColors.textMuted.withOpacity(0.5)),
          border: InputBorder.none,
        ),
      ),
    );
  }

  Widget _buildResultsList() {
    return Expanded(
      child: BlocBuilder<SearchBloc, SearchState>(
        builder: (context, state) {
          if (state is SearchLoading) {
            return _buildSkeleton();
          }
          if (state is SearchError) {
            return Center(child: Text(state.message, style: const TextStyle(color: AppColors.error)));
          }
          if (state is SearchSuccess) {
            if (state.results.isEmpty) {
              return const Center(child: Text("Nenhum resultado encontrado."));
            }
            return ListView.builder(
              itemCount: state.results.length,
              itemBuilder: (context, index) {
                return _buildAnimeCard(state.results[index]);
              },
            );
          }
          return const Center(child: Text("Digite algo para começar..."));
        },
      ),
    );
  }

  Widget _buildAnimeCard(Anime anime) {
    return Container(
      margin: const EdgeInsets.only(bottom: 15),
      height: 120,
      child: InkWell(
        onTap: () {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (context) => DetailsScreen(anime: anime),
            ),
          );
        },
        child: GlassCard(
          padding: EdgeInsets.zero,
          child: Row(
            children: [
               ClipRRect(
                 borderRadius: const BorderRadius.only(topLeft: Radius.circular(16), bottomLeft: Radius.circular(16)),
                 child: CachedNetworkImage(
                   imageUrl: anime.imageUrl,
                   width: 85,
                   height: 120,
                   fit: BoxFit.cover,
                   placeholder: (context, url) => Container(color: Colors.white10),
                   errorWidget: (context, url, error) => const Icon(Icons.error),
                 ),
               ),
               const SizedBox(width: 15),
               Expanded(
                 child: Padding(
                   padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 8),
                   child: Column(
                     crossAxisAlignment: CrossAxisAlignment.start,
                     children: [
                       Text(
                         anime.name,
                         maxLines: 2,
                         overflow: TextOverflow.ellipsis,
                         style: GoogleFonts.outfit(
                           fontSize: 18,
                           fontWeight: FontWeight.bold,
                           color: AppColors.textMain,
                         ),
                       ),
                       const Spacer(),
                       Row(
                         children: [
                           Container(
                             padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                             decoration: BoxDecoration(
                               color: AppColors.primary.withOpacity(0.2),
                               borderRadius: BorderRadius.circular(8),
                             ),
                             child: Text(
                               anime.source,
                               style: const TextStyle(color: AppColors.primary, fontSize: 10, fontWeight: FontWeight.bold),
                             ),
                           ),
                         ],
                       )
                     ],
                   ),
                 ),
               )
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildSkeleton() {
    return Shimmer.fromColors(
      baseColor: Colors.white12,
      highlightColor: Colors.white24,
      child: ListView.builder(
        itemCount: 5,
        itemBuilder: (_, __) => Padding(
          padding: const EdgeInsets.only(bottom: 15),
          child: Container(
            height: 120,
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(16),
            ),
          ),
        ),
      ),
    );
  }
}
