package com.samaya.gotictactoe;

import android.content.DialogInterface;
import android.os.Bundle;
import android.view.Menu;
import android.view.MenuItem;

import androidx.appcompat.app.AlertDialog;
import androidx.appcompat.app.AppCompatActivity;

import com.tictactoe.tictactoe.mobile.EbitenView;
import com.tictactoe.tictactoe.mobile.IGameCallback;
import com.tictactoe.tictactoe.mobile.Mobile;

import go.Seq;

public class MainActivity extends AppCompatActivity implements IGameCallback {

    public static final int HUMAN_PLAYER  = 1;
    public static final int AI_PLAYER = -1;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

           setContentView(R.layout.activity_main);
        Seq.setContext(getApplicationContext());

        Mobile.registerGameCallback(this);

    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        getMenuInflater().inflate(R.menu.appmenu, menu);
        return true;
    }

    private EbitenView getEbitenView() {
        return this.findViewById(R.id.ebitenview);
    }

    @Override
    public void onBackPressed() {
        super.onBackPressed();
        finish();
    }

    @Override
    protected void onPause() {
        super.onPause();
        this.getEbitenView().suspendGame();
    }

    @Override
    protected void onResume() {
        super.onResume();
        this.getEbitenView().resumeGame();
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle item selection
        if(item.getItemId() == R.id.new_game){
            Mobile.playAgain();
            return  true;
        } else if(item.getItemId() == R.id.about) {
            //TODO: show about dialog
            return true;
        }
        return super.onOptionsItemSelected(item);

    }

    @Override
    public void gameOverCallBack(long winner) {
        runOnUiThread(new Runnable() {
            @Override
            public void run() {
                String message  =  winner == HUMAN_PLAYER? "You Win" : winner == AI_PLAYER ? "You Lose" : "Tide the Game";
                AlertDialog.Builder builder = new AlertDialog.Builder(MainActivity.this);
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