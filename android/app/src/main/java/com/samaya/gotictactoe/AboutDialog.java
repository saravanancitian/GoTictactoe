package com.samaya.gotictactoe;

import android.app.AlertDialog;
import android.app.Dialog;
import android.content.DialogInterface;
import android.content.res.AssetManager;
import android.os.Bundle;
import android.text.Html;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.TextView;

import androidx.annotation.NonNull;
import androidx.annotation.Nullable;
import androidx.fragment.app.DialogFragment;

import com.google.android.material.dialog.MaterialAlertDialogBuilder;
import com.tictactoe.tictactoe.mobile.Mobile;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;

public class AboutDialog extends DialogFragment {

    String aboutText = "About Text";
    @Override
    public void onCreate(@Nullable Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        AssetManager assetManager = getActivity().getAssets();
        InputStream is = null;
        try {
            is = assetManager.open("about.html");
            byte data[] = readFully(is);
            aboutText =  data!= null? new String(data) :" About";

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
    }


    public byte[] readFully(InputStream is) throws IOException{
        byte retval[] = null;

        if(is != null){
            ByteArrayOutputStream baos = new ByteArrayOutputStream();
            byte packet[] = new byte[1024];

            while(is.read(packet) > 0 ){
                baos.write(packet);
            }
            retval = baos.toByteArray();
            baos.close();
        }
        return retval;
    }

    // Your own onCreate_Dialog_View method
    public View onCreateDialogView(LayoutInflater inflater, @Nullable ViewGroup container, @Nullable Bundle savedInstanceState) {
        return inflater.inflate(R.layout.about, container); // inflate here
    }

    @Override
    public void onViewCreated(View view, @Nullable Bundle savedInstanceState) {
        super.onViewCreated(view, savedInstanceState);

        TextView txtAbout = view.findViewById(R.id.txt_about);
        txtAbout.setText(Html.fromHtml(aboutText));
    }

    @NonNull
    @Override
    public Dialog onCreateDialog(@Nullable Bundle savedInstanceState) {
        MaterialAlertDialogBuilder builder = new MaterialAlertDialogBuilder(getActivity())

                .setNegativeButton("Cancel",new DialogInterface.OnClickListener(){

                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        Mobile.resume();
                        dialog.dismiss();
                    }
                });

        View view = onCreateDialogView(getActivity().getLayoutInflater(), null, null);
        onViewCreated(view, null);
        builder.setView(view);
        Dialog dialog = builder.create();
        return dialog;
    }
    @Override
    public void onCancel(DialogInterface dialog) {
        Mobile.resume();

    }
}
