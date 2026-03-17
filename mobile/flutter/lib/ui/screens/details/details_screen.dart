import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:cached_network_image/cached_network_image.dart';
import '../../../core/app_colors.dart';
import '../../../models/anime.dart';
import '../../widgets/glass_card.dart';
import 'details_bloc.dart';
import '../player/player_screen.dart';

class DetailsScreen extends StatelessWidget {
  final Anime anime;

  const DetailsScreen({Key? key, required this.anime}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (context) => DetailsBloc()..add(LoadEpisodes(anime.url, anime.source)),
      child: Scaffold(
        backgroundColor: AppColors.background,
        body: CustomScrollView(
          slivers: [
            _buildSliverAppBar(context),
            SliverToBoxAdapter(
              child: Padding(
                padding: const EdgeInsets.all(20),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _buildInfoSection(),
                    const SizedBox(height: 25),
                    _buildDescription(),
                    const SizedBox(height: 30),
                    _buildEpisodesHeader(),
                    const SizedBox(height: 15),
                  ],
                ),
              ),
            ),
            _buildEpisodesList(),
          ],
        ),
      ),
    );
  }

  Widget _buildSliverAppBar(BuildContext context) {
    return SliverAppBar(
      expandedHeight: 350,
      backgroundColor: AppColors.background,
      pinned: true,
      flexibleSpace: FlexibleSpaceBar(
        background: Stack(
          fit: StackFit.expand,
          children: [
            CachedNetworkImage(
              imageUrl: anime.imageUrl,
              fit: BoxFit.cover,
            ),
            Container(
              decoration: const BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topCenter,
                  end: Alignment.bottomCenter,
                  colors: [Colors.black26, AppColors.background],
                ),
              ),
            ),
          ],
        ),
      ),
      leading: IconButton(
        icon: const Icon(Icons.arrow_back_ios_new, color: Colors.white),
        onPressed: () => Navigator.pop(context),
      ),
    );
  }

  Widget _buildInfoSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          anime.name,
          style: GoogleFonts.outfit(
            fontSize: 28,
            fontWeight: FontWeight.bold,
            color: AppColors.textMain,
          ),
        ),
        const SizedBox(height: 8),
        Row(
          children: [
            Text(
              anime.source.toUpperCase(),
              style: const TextStyle(color: AppColors.accent, fontWeight: FontWeight.bold, fontSize: 12),
            ),
            if (anime.details?.status != "") ...[
              const SizedBox(width: 10),
              const CircleAvatar(radius: 2, backgroundColor: Colors.white24),
              const SizedBox(width: 10),
              Text(
                anime.details?.status ?? "",
                style: const TextStyle(color: AppColors.textMuted, fontSize: 12),
              ),
            ],
            const Spacer(),
            if (anime.details?.averageScore != 0)
              Row(
                children: [
                  const Icon(Icons.star, color: Colors.amber, size: 16),
                  const SizedBox(width: 4),
                  Text(
                    "${anime.details?.averageScore}%",
                    style: const TextStyle(color: Colors.amber, fontWeight: FontWeight.bold),
                  ),
                ],
              ),
          ],
        ),
      ],
    );
  }

  Widget _buildDescription() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surface.withOpacity(0.5),
        borderRadius: BorderRadius.circular(16),
      ),
      child: Text(
        anime.details?.description.replaceAll(RegExp(r'<[^>]*>'), '') ?? "Sem descrição disponível.",
        maxLines: 4,
        overflow: TextOverflow.ellipsis,
        style: TextStyle(color: AppColors.textMuted, height: 1.5),
      ),
    );
  }

  Widget _buildEpisodesHeader() {
    return Text(
      "Episódios",
      style: GoogleFonts.outfit(
        fontSize: 22,
        fontWeight: FontWeight.bold,
        color: AppColors.textMain,
      ),
    );
  }

  Widget _buildEpisodesList() {
    return BlocBuilder<DetailsBloc, DetailsState>(
      builder: (context, state) {
        if (state is DetailsLoading) {
          return const SliverFillRemaining(child: Center(child: CircularProgressIndicator()));
        }
        if (state is DetailsError) {
          return SliverToBoxAdapter(child: Center(child: Text(state.message)));
        }
        if (state is DetailsSuccess) {
          return SliverList(
            delegate: SliverChildBuilderDelegate(
              (context, index) {
                final ep = state.episodes[index];
                return Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
                  child: InkWell(
                    onTap: () {
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (context) => PlayerScreen(anime: anime, episode: ep),
                        ),
                      );
                    },
                    child: GlassCard(
                      padding: const EdgeInsets.all(15),
                      child: Row(
                        children: [
                          Container(
                            width: 40,
                            height: 40,
                            decoration: BoxDecoration(
                              color: AppColors.primary.withOpacity(0.2),
                              shape: BoxShape.circle,
                            ),
                            child: Center(
                              child: Text(
                                ep.number,
                                style: const TextStyle(color: AppColors.primary, fontWeight: FontWeight.bold),
                              ),
                            ),
                          ),
                          const SizedBox(width: 15),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  ep.title ?? "Episódio ${ep.number}",
                                  style: const TextStyle(color: AppColors.textMain, fontWeight: FontWeight.w600),
                                ),
                                if (ep.isFiller)
                                  const Text("Filler", style: TextStyle(color: Colors.orange, fontSize: 10)),
                              ],
                            ),
                          ),
                          const Icon(Icons.play_arrow_rounded, color: AppColors.accent),
                        ],
                      ),
                    ),
                  ),
                );
              },
              childCount: state.episodes.length,
            ),
          );
        }
        return const SliverToBoxAdapter(child: SizedBox.shrink());
      },
    );
  }
}
