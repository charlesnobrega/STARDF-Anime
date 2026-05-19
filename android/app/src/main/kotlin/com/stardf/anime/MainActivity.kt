package com.stardf.anime

import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel

class MainActivity: FlutterActivity() {
    private val CHANNEL = "com.stardf.anime/bridge"

    override fun configureFlutterEngine(flutterEngine: FlutterEngine) {
        super.configureFlutterEngine(flutterEngine)

        MethodChannel(flutterEngine.dartExecutor.binaryMessenger, CHANNEL).setMethodCallHandler { call, result ->
            when (call.method) {
                "getAnimes" -> {
                    // TODO: Call Go backend via GoMobile
                    result.success("[]")
                }
                "addToWatchlist" -> {
                    val animeId = call.argument<String>("animeId")
                    // TODO: Call Go backend via GoMobile
                    result.success(true)
                }
                "markAsWatched" -> {
                    val animeId = call.argument<String>("animeId")
                    val episodeId = call.argument<String>("episodeId")
                    // TODO: Call Go backend via GoMobile
                    result.success(true)
                }
                "getSyncStatus" -> {
                    // TODO: Call Go backend via GoMobile
                    result.success("{}")
                }
                "syncWithDesktop" -> {
                    // TODO: Call Go backend via GoMobile
                    result.success(true)
                }
                else -> result.notImplemented()
            }
        }
    }
}
