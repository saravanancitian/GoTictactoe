package com.samaya.gotictactoe;

import android.content.Context;
import android.util.AttributeSet;

import com.tictactoe.tictactoe.mobile.EbitenView;

public class EbitenViewWithErrorHandling extends EbitenView {
    public EbitenViewWithErrorHandling(Context context) {
        super(context);
    }

    public EbitenViewWithErrorHandling(Context context, AttributeSet attributeSet) {
        super(context, attributeSet);
    }

    @Override
    protected void onErrorOnGameUpdate(Exception e) {
        // You can define your own error handling e.g., using Crashlytics.
        // e.g., Crashlytics.logException(e);

        if(getContext() instanceof  MainActivity) {
            ((MainActivity) getContext()).getCrashlytics().recordException(e);
        }
        super.onErrorOnGameUpdate(e);
    }
}
