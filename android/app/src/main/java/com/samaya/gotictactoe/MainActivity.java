package com.samaya.gotictactoe;

import android.content.res.AssetManager;
import android.os.Bundle;
import android.text.Html;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.CompoundButton;
import android.widget.FrameLayout;

import androidx.appcompat.app.AlertDialog;
import androidx.appcompat.app.AppCompatActivity;

import com.google.android.gms.ads.AdRequest;
import com.google.android.gms.ads.AdView;
import com.google.android.material.dialog.MaterialAlertDialogBuilder;
import com.google.android.material.switchmaterial.SwitchMaterial;
import com.google.android.material.textview.MaterialTextView;
import com.google.firebase.analytics.FirebaseAnalytics;
import com.google.firebase.crashlytics.FirebaseCrashlytics;
import com.tictactoe.tictactoe.mobile.EbitenView;
import com.tictactoe.tictactoe.mobile.IGameCallback;
import com.tictactoe.tictactoe.mobile.Mobile;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.ObjectInputStream;
import java.io.ObjectOutputStream;
import java.util.Date;

import go.Seq;

public class MainActivity extends AppCompatActivity implements IGameCallback {

    public static final String SCORE_FILE = "score.dat";
    public static final String ABOUT_FILE = "about.html";
    public static final int HUMAN_PLAYER  = 1;
    public static final int AI_PLAYER = -1;
    public static final int GAME_TIED = 0;

    private AdView adView;

    private  AlertDialog aboutDialog;
    private  AlertDialog resetAlertDialog;
    private AlertDialog scoreDialog;
    private AlertDialog settingDialog;

    private AlertDialog gameOverDialog;

    private ScoreSettings score;
    AssetManager assetManager;

    String aboutString = "";

    View scoreview;
    View settingview;
    MaterialTextView txtGameOver;
    EbitenView ebitenView;


    public FirebaseAnalytics getmFirebaseAnalytics() {
        return mFirebaseAnalytics;
    }

    public FirebaseCrashlytics getCrashlytics() {
        return crashlytics;
    }

    private FirebaseAnalytics mFirebaseAnalytics;
    private FirebaseCrashlytics crashlytics;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        assetManager = getAssets();
        mFirebaseAnalytics = FirebaseAnalytics.getInstance(this);
        crashlytics = FirebaseCrashlytics.getInstance();
        setContentView(R.layout.activity_main);
        Seq.setContext(getApplicationContext());
        Mobile.registerGameCallback(this);
        adView = findViewById(R.id.adView);
        AdRequest adRequest = new AdRequest.Builder().build();
        File file = new File(getFilesDir(), SCORE_FILE);
        if(file.exists()){
            loadScore();
        } else {
            score = ScoreSettings.getInstance();
        }

        adView.loadAd(adRequest);
        aboutString = readAssetFile(ABOUT_FILE);
        createAboutDialog();
        createSettingDialog();
        createScoreDialog();
        createResetAlertDialog();
        createGameOverDialog();
        createEbitenView();
        SwitchMaterial sndSwitch = settingview.findViewById(R.id.snd_switch);
        sndSwitch.setChecked(score.isSettingsSound());

        sndSwitch.setOnCheckedChangeListener(new CompoundButton.OnCheckedChangeListener() {
            @Override
            public void onCheckedChanged(CompoundButton buttonView, boolean isChecked) {
                if (isChecked){
                    Mobile.setSoundOff(false);
                } else {
                    Mobile.setSoundOff(true);
                }
                score.setSettingsSound(isChecked);
            }
        });

        SwitchMaterial timerSwitch = settingview.findViewById(R.id.timer_switch);
        timerSwitch.setChecked(score.isSettingShowTimer());
        timerSwitch.setOnCheckedChangeListener(new CompoundButton.OnCheckedChangeListener() {
            @Override
            public void onCheckedChanged(CompoundButton buttonView, boolean isChecked) {
                if(isChecked){
                    Mobile.setShowTimerOff(false);

                } else {
                    Mobile.setShowTimerOff(true);
                }
                score.setSettingShowTimer(isChecked);
            }
        });
    }


    private void createEbitenView(){
        ebitenView = new EbitenViewWithErrorHandling(this);
        FrameLayout.LayoutParams params = new FrameLayout.LayoutParams(FrameLayout.LayoutParams.MATCH_PARENT, FrameLayout.LayoutParams.MATCH_PARENT);
        ebitenView.setLayoutParams(params);

        FrameLayout gameframe = findViewById(R.id.gameframe);
        gameframe.addView(ebitenView);

    }

    private void loadScore(){
        FileInputStream fis  = null;
        ObjectInputStream ois = null;
       try{
           fis = openFileInput(SCORE_FILE);
           ois = new ObjectInputStream(fis);
           score = (ScoreSettings) ois.readObject();

       } catch (IOException | ClassNotFoundException e) {
           Log.e("Load Score", e.getMessage());
           crashlytics.recordException(e);

           score = ScoreSettings.getInstance();


       } finally {
           if(ois != null){
               try {
                   ois.close();
               } catch (IOException e) {
                   Log.e("Load Score", e.getMessage());

               }
           }
           if(fis != null){
               try {
                   fis.close();
               } catch (IOException e) {
                   Log.e("Load Score", e.getMessage());
               }
           }
       }
    }

    void saveScore(){
        FileOutputStream fos = null;
        ObjectOutputStream oos = null;
        try {
            fos = openFileOutput(SCORE_FILE, MODE_PRIVATE);
            oos = new ObjectOutputStream(fos);
            oos.writeObject(score);
        } catch (IOException e) {
            Log.e("SAVE SCORE", e.getMessage());
            crashlytics.recordException(e);

        } finally {
            if(oos != null){
                try {
                    oos.close();
                } catch (IOException e) {
                    Log.e("SAVE SCORE", e.getMessage());
                }
            }
            if(fos != null){
                try {
                    fos.close();
                } catch (IOException e) {
                    Log.e("SAVE SCORE", e.getMessage());
                }
            }
        }
    }

    String readAssetFile(String filename){
        String ret = null;
        InputStream is = null;
        try {
            is = assetManager.open(filename);
            byte data[] = Util.readFully(is);
            ret = data == null? "" : new String(data);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
        finally {
            if( is != null){
                try {
                    is.close();
                } catch (IOException e) {
                    throw new RuntimeException(e);
                }
            }
        }

        return ret;
    }

    void createResetAlertDialog(){
        MaterialAlertDialogBuilder builder = new MaterialAlertDialogBuilder(this);
        builder.setMessage(R.string.reset_alert_text);
        builder.setPositiveButton(R.string.btn_yes, (dialog, which) -> score.reset());
        builder.setNegativeButton(R.string.btn_no, (dialog, which) -> dialog.cancel());
        resetAlertDialog = builder.create();
    }

    void createGameOverDialog(){

        View gameoverView = getLayoutInflater().inflate(R.layout.game_over, null, false);

        txtGameOver = gameoverView.findViewById(R.id.txt_gameover);

        MaterialAlertDialogBuilder builder = new MaterialAlertDialogBuilder(MainActivity.this);
        builder.setTitle(R.string.game_over);
        builder.setView(gameoverView);
        builder.setPositiveButton(R.string.btn_play_again, (dialog, id) -> Mobile.playAgain(score.getTotalPlayed(), score.getTotalWin()));
        builder.setNegativeButton(R.string.btn_exit, (dialog, id) -> finish());
        builder.setNeutralButton(R.string.btn_view_board, (dialog, which) -> dialog.cancel());

        gameOverDialog = builder.create();
    }

    void createSettingDialog(){

        settingview = getLayoutInflater().inflate(R.layout.settings,null, false);
        MaterialAlertDialogBuilder builder = new MaterialAlertDialogBuilder(this);

        builder.setView(settingview);
        builder.setNegativeButton("Close", (dialog, which) -> dialog.dismiss());
        settingDialog = builder.create();
    }

    void createScoreDialog(){

        scoreview = getLayoutInflater().inflate(R.layout.score,null, false);
        MaterialAlertDialogBuilder builder = new MaterialAlertDialogBuilder(this);

        builder.setView(scoreview);
        builder.setNegativeButton(R.string.btn_cancel, (dialog, which) -> dialog.dismiss());
        builder.setNeutralButton(R.string.btn_reset, (dialog, which) -> resetAlertDialog.show());
        scoreDialog = builder.create();
    }

    void createAboutDialog(){
        View aboutView = getLayoutInflater().inflate(R.layout.about, null, false);
        MaterialAlertDialogBuilder builder = new MaterialAlertDialogBuilder(this);
        builder.setView(aboutView);
        builder.setNegativeButton(R.string.btn_cancel, (dialog, which) -> dialog.dismiss());
        MaterialTextView txt_about = aboutView.findViewById(R.id.txt_about);
        txt_about.setText(Html.fromHtml(aboutString));
        aboutDialog = builder.create();
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        getMenuInflater().inflate(R.menu.appmenu, menu);
        return true;
    }

    @Override
    public void onBackPressed() {
        super.onBackPressed();
        finish();;
    }

    @Override
    protected void onPause() {
        if (adView != null) {
            adView.pause();
        }
        saveScore();
        pauseGame();
        super.onPause();
    }

    @Override
    protected void onResume() {
        super.onResume();
        if (adView != null) {
            adView.resume();
        }
        resumeGame();
    }

    @Override
    public void onDestroy() {
        if (adView != null) {
            adView.destroy();
        }
        Mobile.destroy();
        super.onDestroy();
        System.exit(0);
    }
    
    void pauseGame(){
        Mobile.pause();
        this.ebitenView.suspendGame();
    }
    
    void resumeGame(){

        Mobile.resume();
        this.ebitenView.resumeGame();
        Mobile.setSoundOff(!score.isSettingsSound());
        Mobile.setShowTimerOff(!score.isSettingShowTimer());
    }

    void showScore(){
        MaterialTextView txtTotalplayed = scoreview.findViewById(R.id.totalPlayed);
        txtTotalplayed.setText(String.valueOf(score.getTotalPlayed()));

        MaterialTextView txtTotalWin = scoreview.findViewById(R.id.totalWin);
        txtTotalWin.setText(String.valueOf(score.getTotalWin()));

        MaterialTextView txtTotalTied = scoreview.findViewById(R.id.totalTied);
        txtTotalTied.setText(String.valueOf(score.getTotalTied()));

        Date dt1 = score.getDate1() > 0? new Date( score.getDate1() ): null;
        long dur1 =  score.getTopPlayedTime1() > 0? score.getTopPlayedTime1()/ ScoreSettings.SEC_IN_MILLIS : 0;
        Date dt2 =  score.getDate2() > 0? new Date(score.getDate2()) : null;
        long dur2 = score.getTopPlayedTime2() > 0? score.getTopPlayedTime2()/ ScoreSettings.SEC_IN_MILLIS : 0;
        Date dt3 =  score.getDate3() > 0?new Date(score.getDate3()) :  null;
        long dur3 =  score.getTopPlayedTime3() > 0? score.getTopPlayedTime3()/ ScoreSettings.SEC_IN_MILLIS : 0;


        MaterialTextView t1 = scoreview.findViewById(R.id.tp1);
        t1.setText(Util.formatTime(dur1));

        MaterialTextView d1 = scoreview.findViewById(R.id.date1);
        d1.setText(Util.formatDate(dt1));

        MaterialTextView t2 = scoreview.findViewById(R.id.tp2);
        t2.setText(Util.formatTime(dur2));

        MaterialTextView d2 = scoreview.findViewById(R.id.date2);
        d2.setText(Util.formatDate(dt2));

        MaterialTextView t3 = scoreview.findViewById(R.id.tp3);
        t3.setText(Util.formatTime(dur3));

        MaterialTextView d3 = scoreview.findViewById(R.id.date3);
        d3.setText(Util.formatDate(dt3));

        scoreDialog.show();
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle item selection
        if(item.getItemId() == R.id.new_game){
            Mobile.playAgain(score.getTotalPlayed(), score.getTotalWin());
            return  true;
        } else if(item.getItemId() == R.id.settings) {
            Bundle bundle = new Bundle();
            bundle.putString(FirebaseAnalytics.Param.SCREEN_NAME, "Settings Dialog");
            bundle.putString(FirebaseAnalytics.Param.SCREEN_CLASS, "MainActivity");
            mFirebaseAnalytics.logEvent(FirebaseAnalytics.Event.SCREEN_VIEW, bundle);

            settingDialog.show();
            return true;
        } else if(item.getItemId() == R.id.score){
            Bundle bundle = new Bundle();
            bundle.putString(FirebaseAnalytics.Param.SCREEN_NAME, "Score Dialog");
            bundle.putString(FirebaseAnalytics.Param.SCREEN_CLASS, "MainActivity");
            mFirebaseAnalytics.logEvent(FirebaseAnalytics.Event.SCREEN_VIEW, bundle);

            showScore();
            return true;
        } else if(item.getItemId() == R.id.about) {

            Bundle bundle = new Bundle();
            bundle.putString(FirebaseAnalytics.Param.SCREEN_NAME, "About Dialog");
            bundle.putString(FirebaseAnalytics.Param.SCREEN_CLASS, "MainActivity");
            mFirebaseAnalytics.logEvent(FirebaseAnalytics.Event.SCREEN_VIEW, bundle);

            aboutDialog.show();
            return true;
        }
        return super.onOptionsItemSelected(item);

    }

    @Override
    public void gameOverCallBack(long winner, long duration) {
        score.addPlayed((int)winner, duration);

        runOnUiThread(() -> {
            try {
                String message  = getString( winner == HUMAN_PLAYER? R.string.you_win : winner == AI_PLAYER ? R.string.you_lose : R.string.tied_game);
                txtGameOver.setText(message);
                gameOverDialog.show();
            } catch (Exception e){
                Log.e("GameOverCallback -runOnUIThread ", e.getMessage());
                MainActivity.this.crashlytics.recordException(e);
            }
        });
    }
}