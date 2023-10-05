package com.samaya.gotictactoe;

import static android.view.ViewGroup.LayoutParams.MATCH_PARENT;
import static android.view.ViewGroup.LayoutParams.WRAP_CONTENT;

import android.content.DialogInterface;
import android.content.res.AssetManager;
import android.os.Bundle;
import android.text.Html;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.webkit.WebView;

import androidx.appcompat.app.AlertDialog;
import androidx.appcompat.app.AppCompatActivity;

import com.google.android.gms.ads.AdRequest;
import com.google.android.gms.ads.AdView;
import com.google.android.material.dialog.MaterialAlertDialogBuilder;
import com.google.android.material.textview.MaterialTextView;
import com.tictactoe.tictactoe.mobile.EbitenView;
import com.tictactoe.tictactoe.mobile.IGameCallback;
import com.tictactoe.tictactoe.mobile.Mobile;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.ObjectInputStream;
import java.io.ObjectOutputStream;

import go.Seq;

public class MainActivity extends AppCompatActivity implements IGameCallback {

    public static final String SCOREFILE = "score.dat";
    public static final int HUMAN_PLAYER  = 1;
    public static final int AI_PLAYER = -1;
    public static final int GAME_TIED = 0;

    private AdView adView;

    private  AlertDialog aboutDialog;
    private  AlertDialog resetAlertDialog;

    private AlertDialog scoreDialog;

    private  Score score;
    AssetManager assetManager;

    String scoreString = "";
    String aboutString = "";

     WebView scoretxtview;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        assetManager = getAssets();
        setContentView(R.layout.activity_main);
        Seq.setContext(getApplicationContext());
        Mobile.registerGameCallback(this);
        adView = findViewById(R.id.adView);
        AdRequest adRequest = new AdRequest.Builder().build();
        File file = new File(getFilesDir(), SCOREFILE);
        if(file.exists()){
            loadScore();
        } else {
            score = Score.getInstance();
        }

        adView.loadAd(adRequest);
        scoreString = readAssetFile("score.html");
        aboutString = readAssetFile("about.html");
        createAboutDialog();
        createScoreDialog();
        createResetAlertDialog();
    }

    private void loadScore(){
        FileInputStream fis  = null;
        ObjectInputStream ois = null;
       try{
           fis = openFileInput(SCOREFILE);
           ois = new ObjectInputStream(fis);
           score = (Score) ois.readObject();

       } catch (IOException | ClassNotFoundException e) {
           Log.e("Load Score", e.getMessage());
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
            fos = openFileOutput(SCOREFILE, MODE_PRIVATE);
            oos = new ObjectOutputStream(fos);
            oos.writeObject(score);
        } catch (FileNotFoundException e) {
            Log.e("SAVE SCORE", e.getMessage());
        } catch (IOException e) {
            Log.e("SAVE SCORE", e.getMessage());
        }
        finally {
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
        builder.setMessage("Resetting the score cannot be undone. Do you wish to reset");
        builder.setPositiveButton("Yes", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                score.reset();
            }
        });
        builder.setNegativeButton("No", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                dialog.cancel();
            }
        });
        resetAlertDialog = builder.create();
    }

    void createScoreDialog(){
        scoretxtview = new WebView(this);

        MaterialAlertDialogBuilder builder = new MaterialAlertDialogBuilder(this);

        builder.setView(scoretxtview);
        builder.setNegativeButton("Cancel", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                dialog.dismiss();
            }
        });
        builder.setNeutralButton("Reset", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                resetAlertDialog.show();

            }
        });
        scoreDialog = builder.create();
    }

    void createAboutDialog(){
        WebView wv = new WebView(this);
        MaterialAlertDialogBuilder builder = new MaterialAlertDialogBuilder(this);
        builder.setView(wv);
        builder.setNegativeButton("Cancel", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                dialog.dismiss();
            }
        });
        wv.loadUrl("file:///android_asset/about.html");
        aboutDialog = builder.create();
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        getMenuInflater().inflate(R.menu.appmenu, menu);
        return true;
    }

    public EbitenView getEbitenView() {
        return this.findViewById(R.id.ebitenview);
    }

    @Override
    public void onBackPressed() {
        super.onBackPressed();
        finish();
    }

    @Override
    protected void onPause() {
        if (adView != null) {
            adView.pause();
        }
        saveScore();
        this.getEbitenView().suspendGame();
        super.onPause();
        //Mobile.resume();
    }

    @Override
    protected void onResume() {
        super.onResume();
        if (adView != null) {
            adView.resume();
        }
        this.getEbitenView().resumeGame();
        //Mobile.pause();

    }

    @Override
    public void onDestroy() {
        if (adView != null) {
            adView.destroy();
        }
        super.onDestroy();
    }

    void showScore(){
        String scoretxt = score.formattedString( scoreString);

        scoretxtview.loadData(scoretxt, "text/html; charset=UTF-8", null);

        scoreDialog.show();
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle item selection
        if(item.getItemId() == R.id.new_game){
            Mobile.playAgain();
            return  true;
        } else if(item.getItemId() == R.id.score){
            showScore();
            return true;
        } else if(item.getItemId() == R.id.about) {
            aboutDialog.show();
            return true;
        }
        return super.onOptionsItemSelected(item);

    }

    @Override
    public void gameOverCallBack(long winner, long duration) {
        score.addPlayed((int)winner, duration);

        runOnUiThread(new Runnable() {
            @Override
            public void run() {
                String message  =  winner == HUMAN_PLAYER? "You Win" : winner == AI_PLAYER ? "You Lose" : "Tide the Game";

                MaterialAlertDialogBuilder builder = new MaterialAlertDialogBuilder(MainActivity.this, R.style.AlertDialogGameover);
                builder.setMessage(message).setTitle("Game Over");

                builder.setPositiveButton("Play Again", new DialogInterface.OnClickListener() {
                    public void onClick(DialogInterface dialog, int id) {
                        Mobile.playAgain();
                    }
                });
                builder.setNegativeButton("Exit", new DialogInterface.OnClickListener() {
                    public void onClick(DialogInterface dialog, int id) {
                        finish();
                    }
                });
                builder.setNeutralButton("View Board", new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        dialog.cancel();
                    }
                });

                AlertDialog dialog = builder.create();
                dialog.show();
            }
        });
    }
}