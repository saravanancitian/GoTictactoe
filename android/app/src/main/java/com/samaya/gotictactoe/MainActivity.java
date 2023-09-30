package com.samaya.gotictactoe;

import android.os.Bundle;
import android.view.View;
import android.widget.Toast;

import androidx.appcompat.app.AppCompatActivity;

import com.tictactoe.tictactoe.mobile.EbitenView;
import com.tictactoe.tictactoe.mobile.IGameCallback;
import com.tictactoe.tictactoe.mobile.Mobile;

import go.Seq;

public class MainActivity extends AppCompatActivity implements IGameCallback {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

           setContentView(R.layout.activity_main);
        Seq.setContext(getApplicationContext());

        Mobile.registerGameCallback(this);

    }

    private EbitenView getEbitenView() {
        return this.findViewById(R.id.ebitenview);
    }

    private View getGameoverLayour() {
        return this.findViewById(R.id.layoutGameover);
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
    public void gameOverCallBack() {
        runOnUiThread(new Runnable() {
            @Override
            public void run() {
                getGameoverLayour().setVisibility(View.VISIBLE);
                Toast.makeText(MainActivity.this, "Game over", Toast.LENGTH_LONG).show();
            }
        });
    }
}