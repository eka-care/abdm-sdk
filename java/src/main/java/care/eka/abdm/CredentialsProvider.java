package care.eka.abdm;

import java.io.IOException;

/**
 * A source of authentication credentials, implementations of which handle
 * automatic refreshing when possible. Mirrors {@code auth.CredentialsProvider}
 * in Go.
 */
public interface CredentialsProvider {

    /** Retrieves the current credentials, refreshing or logging in as needed. */
    Credentials retrieve() throws IOException, InterruptedException;
}
